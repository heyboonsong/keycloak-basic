import Keycloak from 'keycloak-js';

const keycloakConfig = {
  url: '',
  realm: '',
  clientId: ''
};

const keycloak = new Keycloak(keycloakConfig);

export default keycloak;