import React from "react";
import { KeycloakProvider } from "./contexts/KeycloakContext";
import { useKeycloak } from "./contexts/KeycloakContext";
import Login from "./components/Login";
import TodoApp from "./components/TodoApp";
import "./App.css";

const AppContent: React.FC = () => {
  const { authenticated } = useKeycloak();

  if (!authenticated) {
    return <Login />;
  }

  return <TodoApp />;
};

const App: React.FC = () => {
  return (
    <KeycloakProvider>
      <AppContent />
    </KeycloakProvider>
  );
};

export default App;
