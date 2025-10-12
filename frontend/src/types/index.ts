export interface Transaction {
  id: number;
  type: 'income' | 'expense';
  amount: number;
  category_id: number;
  category?: Category;
  transaction_date: string;
  memo: string;
  created_at: string;
  updated_at: string;
}

export interface Category {
  id: number;
  name: string;
  type: 'income' | 'expense';
  color: string;
  created_at: string;
  updated_at: string;
}

export interface Budget {
  id: number;
  category_id: number;
  category?: Category;
  amount: number;
  target_year: number;
  target_month: number;
  created_at: string;
  updated_at: string;
}

export interface CategorySummary {
  category_id: number;
  category_name: string;
  category_type: string;
  total: number;
  budget: number;
  percentage: number;
}

export interface MonthlySummary {
  year: number;
  month: number;
  total_income: number;
  total_expense: number;
  balance: number;
  category_summary: Record<number, CategorySummary>;
}

export interface CreateTransactionRequest {
  type: 'income' | 'expense';
  amount: number;
  category_id: number;
  transaction_date: string;
  memo?: string;
}

export interface CreateCategoryRequest {
  name: string;
  type: 'income' | 'expense';
  color?: string;
}

export interface CreateBudgetRequest {
  category_id: number;
  amount: number;
  target_year: number;
  target_month: number;
}

// API error types for better error handling
export interface ApiError {
  message: string;
  statusCode?: number;
  originalError?: unknown;
  timestamp: string;
}

export class ApplicationError extends Error {
  constructor(
    message: string,
    public statusCode?: number,
    public originalError?: unknown
  ) {
    super(message);
    this.name = 'ApplicationError';
    // Maintain proper stack trace for where our error was thrown
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ApplicationError);
    }
  }

  toJSON(): ApiError {
    return {
      message: this.message,
      statusCode: this.statusCode,
      originalError: this.originalError,
      timestamp: new Date().toISOString(),
    };
  }
}
