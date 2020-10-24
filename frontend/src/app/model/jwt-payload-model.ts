export interface JwtPayloadModel {
  sub: string;
  iat: number;
  exp: number;
  name: string;
}
