package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

// User 定义了用户信息的结构体
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // 注意：实际开发中密码应当加密存储
}

// usersDB 代表了简单的用户数据库
// 注意：这是一个简化的示例，不适合生产环境
var usersDB = make(map[string]User)
var dbMutex = &sync.Mutex{} // 使用互斥锁确保并发安全

// registerHandler 处理用户注册
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	if _, exists := usersDB[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	usersDB[user.Username] = user
	w.WriteHeader(http.StatusCreated)
}

// loginHandler 处理用户登录
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbMutex.Lock()
	defer dbMutex.Unlock()

	storedUser, exists := usersDB[user.Username]
	if !exists || storedUser.Password != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

func main() {
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)

	// 启动HTTP服务器
	http.ListenAndServe(":8080", nil)
}
