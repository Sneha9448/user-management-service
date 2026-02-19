import { gql } from 'urql';

export const GET_USERS_QUERY = gql`
  query GetUsers {
    users {
      id
      name
      email
      role
    }
  }
`;

export const GET_ME_QUERY = gql`
  query GetMe {
    me {
      id
      name
      email
    }
  }
`;
