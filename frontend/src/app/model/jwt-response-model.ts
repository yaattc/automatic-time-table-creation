import { User } from './user';

export interface JwtResponseModel {
  token: string;
  user: User;
}
