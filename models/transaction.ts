interface Transaction {
  id: string;
  amount: number;
  description: string;
  category: string;
  date: Date;
  userId: string;
  isIncome: boolean;
}
