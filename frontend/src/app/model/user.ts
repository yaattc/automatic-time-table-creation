export interface User {
  name: string;
  id: string;
  picture: string;
  email: string;
  attrs: {
    privileges: string[];
  };
}
