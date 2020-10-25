export interface ErrorResponseModel {
  timestamp: string;
  status: number;
  error: string;
  message: string;
  path: string;
}
