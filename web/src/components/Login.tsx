import React from "react";
import { useKeycloak } from "../contexts/KeycloakContext";

const Login: React.FC = () => {
  const { login } = useKeycloak();

  return (
    <div className="login-container">
      <div className="login-card">
        <h1>Todo App</h1>
        <p>Please login to manage your todos</p>
        <button onClick={login} className="login-button">
          Login with Keycloak
        </button>
      </div>
    </div>
  );
};

export default Login;
