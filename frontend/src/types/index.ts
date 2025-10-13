/**
 * 取引データの型定義
 */
export interface Transaction {
  /** 取引ID */
  id: number;
  /** 取引種別（収入/支出） */
  type: 'income' | 'expense';
  /** 金額 */
  amount: number;
  /** カテゴリID */
  category_id: number;
  /** カテゴリ情報（結合時に含まれる） */
  category?: Category;
  /** 取引日 */
  transaction_date: string;
  /** メモ */
  memo: string;
  /** 作成日時 */
  created_at: string;
  /** 更新日時 */
  updated_at: string;
}

/**
 * カテゴリデータの型定義
 */
export interface Category {
  /** カテゴリID */
  id: number;
  /** カテゴリ名 */
  name: string;
  /** カテゴリ種別（収入/支出） */
  type: 'income' | 'expense';
  /** 表示色（16進数カラーコード） */
  color: string;
  /** 作成日時 */
  created_at: string;
  /** 更新日時 */
  updated_at: string;
}

/**
 * 予算データの型定義
 */
export interface Budget {
  /** 予算ID */
  id: number;
  /** カテゴリID */
  category_id: number;
  /** カテゴリ情報（結合時に含まれる） */
  category?: Category;
  /** 予算金額 */
  amount: number;
  /** 対象年 */
  target_year: number;
  /** 対象月（1-12） */
  target_month: number;
  /** 作成日時 */
  created_at: string;
  /** 更新日時 */
  updated_at: string;
}

/**
 * カテゴリ別集計データの型定義
 */
export interface CategorySummary {
  /** カテゴリID */
  category_id: number;
  /** カテゴリ名 */
  category_name: string;
  /** カテゴリ種別 */
  category_type: string;
  /** 合計金額 */
  total: number;
  /** 予算金額 */
  budget: number;
  /** 予算に対する使用率（%） */
  percentage: number;
}

/**
 * 月次サマリーデータの型定義
 */
export interface MonthlySummary {
  /** 年 */
  year: number;
  /** 月（1-12） */
  month: number;
  /** 総収入 */
  total_income: number;
  /** 総支出 */
  total_expense: number;
  /** 残高（収入 - 支出） */
  balance: number;
  /** カテゴリ別集計（キー: カテゴリID） */
  category_summary: Record<number, CategorySummary>;
}

/**
 * 取引作成リクエストの型定義
 */
export interface CreateTransactionRequest {
  /** 取引種別（収入/支出） */
  type: 'income' | 'expense';
  /** 金額 */
  amount: number;
  /** カテゴリID */
  category_id: number;
  /** 取引日（YYYY-MM-DD形式） */
  transaction_date: string;
  /** メモ（任意） */
  memo?: string;
}

/**
 * カテゴリ作成リクエストの型定義
 */
export interface CreateCategoryRequest {
  /** カテゴリ名 */
  name: string;
  /** カテゴリ種別（収入/支出） */
  type: 'income' | 'expense';
  /** 表示色（16進数カラーコード、任意） */
  color?: string;
}

/**
 * 予算作成リクエストの型定義
 */
export interface CreateBudgetRequest {
  /** カテゴリID */
  category_id: number;
  /** 予算金額 */
  amount: number;
  /** 対象年 */
  target_year: number;
  /** 対象月（1-12） */
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
