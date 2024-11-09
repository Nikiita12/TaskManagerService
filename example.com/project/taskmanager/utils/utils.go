package utils

import (
    "errors"
    "strings"
)

// UserInfo содержит информацию о пользователе
type UserInfo struct {
    Username  string
    Roles     []string
    Permissions []string
}

// ValidateAuthToken проверяет токен и возвращает информацию о пользователе
func ValidateAuthToken(token string) (UserInfo, error) {
    // Пример проверки токена. На практике здесь должна быть ваша логика проверки.
    if token == "" {
        return UserInfo{}, errors.New("invalid token")
    }

    // Пример декодирования токена. На практике это будет более сложная логика.
    parts := strings.Split(token, ":")
    if len(parts) != 2 {
        return UserInfo{}, errors.New("invalid token format")
    }

    username := parts[0]
    roles := strings.Split(parts[1], ",")
    
    // Пример возвращаемой информации о пользователе
    return UserInfo{
        Username:  username,
        Roles:     roles,
        Permissions: []string{"create_task", "update_task", "view_task"}, // Укажите актуальные разрешения
    }, nil
}

// HasPermission проверяет, есть ли у пользователя необходимое разрешение
func (u UserInfo) HasPermission(permission string) bool {
    for _, p := range u.Permissions {
        if p == permission {
            return true
        }
    }
    return false
}
