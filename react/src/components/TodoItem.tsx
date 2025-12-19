import React from "react";
import type { Todo } from "../types/todo";

interface TodoItemProps {
  todo: Todo;
}

const TodoItem: React.FC<TodoItemProps> = ({ todo }) => {
  return (
    <div className={`todo-item ${todo.completed ? "completed" : ""}`}>
      <span className="todo-text">{todo.text}</span>
    </div>
  );
};

export default TodoItem;
