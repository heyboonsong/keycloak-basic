import React, { createContext, useContext, useEffect, useState } from "react";
import Keycloak from "keycloak-js";
import keycloak from "../config/keycloak";

interface KeycloakContextType {
  keycloak: Keycloak | null;
  authenticated: boolean;
  login: () => void;
  logout: () => void;
}

const KeycloakContext = createContext<KeycloakContextType | undefined>(
  undefined
);

export const useKeycloak = () => {
  const context = useContext(KeycloakContext);
  if (context === undefined) {
    throw new Error("useKeycloak must be used within a KeycloakProvider");
  }
  return context;
};

interface KeycloakProviderProps {
  children: React.ReactNode;
}

export const KeycloakProvider: React.FC<KeycloakProviderProps> = ({
  children,
}) => {
  const [keycloakInstance, setKeycloakInstance] = useState<Keycloak | null>(
    null
  );
  const [authenticated, setAuthenticated] = useState(false);

  useEffect(() => {
    keycloak
      .init({ onLoad: "check-sso" })
      .then((authenticated) => {
        setKeycloakInstance(keycloak);
        setAuthenticated(authenticated);
      })
      .catch((error) => {
        console.error("Keycloak initialization failed:", error);
      });
  }, []);

  const login = () => {
    if (keycloakInstance) {
      keycloakInstance.login();
    }
  };

  const logout = () => {
    if (keycloakInstance) {
      keycloakInstance.logout();
    }
  };

  return (
    <KeycloakContext.Provider
      value={{
        keycloak: keycloakInstance,
        authenticated,
        login,
        logout,
      }}
    >
      {children}
    </KeycloakContext.Provider>
  );
};
