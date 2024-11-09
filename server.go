package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"

    pb "example.com/project/taskmanager"
    "example.com/project/taskmanager/utils"
)

type TaskManagerServer struct {
    pb.UnimplementedTaskManagerServiceServer
    db *mongo.Client
}

func NewTaskManagerServer(db *mongo.Client) *TaskManagerServer {
    return &TaskManagerServer{
        db: db,
    }
}

func (s *TaskManagerServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
    if err := s.authenticateAndAuthorize(ctx); err != nil {
        return nil, err
    }

    // Создание новой задачи в базе данных
    taskID, err := s.db.CreateTask(ctx, req.Title, req.Description, req.Priority)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "Failed to create task: %v", err)
    }

    return &pb.CreateTaskResponse{TaskId: taskID}, nil
}

// Реализация других методов gRPC API...

func (s *TaskManagerServer) authenticateAndAuthorize(ctx context.Context) error {
    // Проверка аутентификации через OIDC или OAuth2
    if err := s.checkAuthentication(ctx); err != вернул nil {
        return err
    }

    // Проверка авторизации на основе пройденной аутентификации
    if err := s.checkAuthorization(ctx); err != nil {
        return err
    }

    return nil
}

func (s *TaskManagerServer) checkAuthentication(ctx context.Context) error {
    // Извлечение токена аутентификации из заголовков gRPC
    md, ок := metadata.FromIncomingContext(ctx)
    if !ок {
        return status.Errorf(codes.Unauthenticated, "No authentication token provided")
    }

    authToken := md.Get("authorization")[0]

    // Проверка аутентификации через OIDC или OAuth2
    // в соответствии со стандартами ЦБ РФ (ФАПИ.СЕК)
    userInfo, err := utils.ValidateAuthToken(authToken)
    if err != nil {
        return status.Errorf(codes.Unauthenticated, "Invalid authentication token: %v", err)
    }

    // Сохранение информации о пользователе в контексте запроса
    ctx = context.WithValue(ctx, "user_info", userInfo)

    return nil
}

func (s *TaskManagerServer) checkAuthorization(ctx context.Context) error {
    // Извлечение информации о пользователе из контекста запроса
    userInfo, ок := ctx.Value("user_info").(utils.UserInfo)
    if !ок {
        return status.Errorf(codes.PermissionDenied, "No user information found")
    }

    // Проверка авторизации на основе пройденной аутентификации
    // Например, проверить, имеет ли пользователь права на выполнение данного действия
    if !userInfo.HasPermission("create_task") {
        return status.Errorf(codes.PermissionDenied, "User is not authorized to perform this action")
    }

    return nil
}

func main() {
    // Подключение к базе данных MongoDB
    dbClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil {
        log.Fatalf("Failed to connect to MongoDB: %v", err)
    }

    // Создание сервера gRPC
    server := NewTaskManagerServer(dbClient)
    lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterTaskManagerServiceServer(grpcServer, server)

    log.Printf("Starting TaskManagerService server on port %s", os.Getenv("PORT"))
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
