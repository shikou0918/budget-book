export const formatNumber = (num: number): string => {
  return new Intl.NumberFormat('ja-JP').format(num);
};

export const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('ja-JP');
};

export const formatCurrency = (amount: number, type: 'income' | 'expense'): string => {
  const sign = type === 'income' ? '+' : '-';
  return `${sign}¥${formatNumber(amount)}`;
};

export const formatTransactionType = (type: 'income' | 'expense'): string => {
  return type === 'income' ? '収入' : '支出';
};
