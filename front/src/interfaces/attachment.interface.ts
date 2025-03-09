export interface Attachment {
  id: number;
  type: "url" | "file";
  url: string;
  filename: string;
  size: number;
}
