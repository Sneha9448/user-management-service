import { gql } from 'urql';

export const REQUEST_OTP_MUTATION = gql`
  mutation RequestOtp($email: String!) {
    requestOtp(email: $email)
  }
`;

export const VERIFY_OTP_MUTATION = gql`
  mutation VerifyOtp($email: String!, $otp: String!, $role: String) {
    verifyOtp(email: $email, otp: $otp, role: $role) {
      token
      user {
        id
        name
        email
        role
      }
    }
  }
`;

export const CREATE_USER_MUTATION = gql`
  mutation CreateUser($name: String!, $email: String!) {
    createUser(name: $name, email: $email) {
      id
      name
      email
    }
  }
`;

export const UPDATE_USER_MUTATION = gql`
  mutation UpdateUser($id: ID!, $name: String!, $email: String!) {
    updateUser(id: $id, name: $name, email: $email) {
      id
      name
      email
    }
  }
`;

export const DELETE_USER_MUTATION = gql`
  mutation DeleteUser($id: ID!) {
    deleteUser(id: $id)
  }
`;
