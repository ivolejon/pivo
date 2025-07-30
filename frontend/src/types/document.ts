export type Document = {
  id: string;
  title: string;
  projectId: string;
  createdAt: string;
  updatedAt: string;
  content: string; // Optional field for document content
  filename: string; // Optional field for the file name
};
