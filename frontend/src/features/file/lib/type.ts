export interface Document {
  id: string;
  data?: {
    title: string;
    documentUrl: string;
  };
  createdAt: string;
  modifiedAt: string;
}

export interface Pageable<D> {
  size: number;
  limit: number;
  total: number;
  data: D[];
  cursor?: string;
}

