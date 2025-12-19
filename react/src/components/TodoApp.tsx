import React, { useState, useEffect } from "react";
import { useKeycloak } from "../contexts/KeycloakContext";
import type { Todo } from "../types/todo";
import TodoList from "./TodoList";

const TodoApp: React.FC = () => {
  const { logout, userName } = useKeycloak();
  const [todos, setTodos] = useState<Todo[]>([]);

  useEffect(() => {
    // Mock data for logged-in users
    const mockTodos: Todo[] = [
      {
        id: "1",
        text: "Learn React with TypeScript",
        completed: true,
        createdAt: new Date(Date.now() - 86400000), // 1 day ago
      },
      {
        id: "2",
        text: "Set up Keycloak SSO",
        completed: true,
        createdAt: new Date(Date.now() - 43200000), // 12 hours ago
      },
      {
        id: "3",
        text: "Style the application with CSS",
        completed: false,
        createdAt: new Date(),
      },
      {
        id: "4",
        text: "Deploy the application",
        completed: false,
        createdAt: new Date(),
      },
    ];

    // Set mock todos when component mounts
    setTodos(mockTodos);
  }, []);

  const completedCount = todos.filter((todo) => todo.completed).length;
  const totalCount = todos.length;

  return (
    <div className="todo-app">
      <header className="todo-header">
        <div className="header-content">
          <h1>My Todos</h1>
          <div className="user-info">
            <span>Welcome, {userName}!</span>
            <button onClick={logout} className="logout-button">
              Logout
            </button>
          </div>
        </div>
        <div className="todo-stats">
          <span>
            {completedCount} of {totalCount} completed
          </span>
        </div>
      </header>

      <main className="todo-main">
        <TodoList todos={todos} />
      </main>
    </div>
  );
};

export default TodoApp;
