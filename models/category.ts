interface Category {
  id: string;
  name: string;
  type: 'income' | 'expense';
  userId: string;
  color?: string;
}
