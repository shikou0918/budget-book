import { describe, test, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import TransactionTable from '../TransactionTable.vue';
import type { Transaction } from '@/types';

const mockTransactions: Transaction[] = [
  {
    id: 1,
    transaction_date: '2024-01-15',
    type: 'income',
    amount: 50000,
    memo: '給与',
    category: {
      id: 1,
      name: '給与',
      type: 'income',
      color: '#4CAF50',
      created_at: '2024-01-15T00:00:00Z',
      updated_at: '2024-01-15T00:00:00Z',
    },
    category_id: 1,
    created_at: '2024-01-15T00:00:00Z',
    updated_at: '2024-01-15T00:00:00Z',
  },
  {
    id: 2,
    transaction_date: '2024-01-16',
    type: 'expense',
    amount: 1200,
    memo: 'ランチ',
    category: {
      id: 2,
      name: '食費',
      type: 'expense',
      color: '#F44336',
      created_at: '2024-01-15T00:00:00Z',
      updated_at: '2024-01-15T00:00:00Z',
    },
    category_id: 2,
    created_at: '2024-01-16T00:00:00Z',
    updated_at: '2024-01-16T00:00:00Z',
  },
];

describe('TransactionTable', () => {
  test('取引データが正しく表示される', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        loading: false,
      },
    });

    expect(wrapper.text()).toContain('2024/1/15');
    expect(wrapper.text()).toContain('給与');
    expect(wrapper.text()).toContain('¥50,000');
    expect(wrapper.text()).toContain('ランチ');
    expect(wrapper.text()).toContain('¥1,200');
  });

  test('loading propがtrueの時にローディング状態を表示する', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: [],
        loading: true,
      },
    });

    // v-data-tableのローディング状態をチェック
    const dataTable = wrapper.findComponent({ name: 'VDataTable' });
    expect(dataTable.props('loading')).toBe(true);
  });

  test('showActionsがtrueの時にアクションボタンを表示する', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showActions: true,
      },
    });

    const editButtons = wrapper.findAll('button').filter(btn => btn.text().includes('編集'));
    const deleteButtons = wrapper.findAll('button').filter(btn => btn.text().includes('削除'));

    expect(editButtons.length).toBe(mockTransactions.length);
    expect(deleteButtons.length).toBe(mockTransactions.length);
  });

  test('showActionsがfalseの時にアクションボタンを非表示にする', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showActions: false,
      },
    });

    const editButtons = wrapper.findAll('button').filter(btn => btn.text().includes('編集'));
    const deleteButtons = wrapper.findAll('button').filter(btn => btn.text().includes('削除'));

    expect(editButtons.length).toBe(0);
    expect(deleteButtons.length).toBe(0);
  });

  test('編集ボタンをクリックした時にeditイベントを発行する', async () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showActions: true,
      },
    });

    const editButton = wrapper.findAll('button').find(btn => btn.text().includes('編集'));

    await editButton?.trigger('click');

    expect(wrapper.emitted('edit')).toBeTruthy();
    expect(wrapper.emitted('edit')?.[0]).toEqual([mockTransactions[0]]);
  });

  test('削除ボタンをクリックした時にdeleteイベントを発行する', async () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showActions: true,
      },
    });

    const deleteButton = wrapper.findAll('button').find(btn => btn.text().includes('削除'));

    await deleteButton?.trigger('click');

    expect(wrapper.emitted('delete')).toBeTruthy();
    expect(wrapper.emitted('delete')?.[0]).toEqual([mockTransactions[0].id]);
  });

  test('showSearchがtrueの時に検索フィールドを表示する', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showSearch: true,
      },
    });

    const searchField = wrapper.findComponent({ name: 'VTextField' });
    expect(searchField.exists()).toBe(true);
    expect(searchField.props('label')).toBe('検索');
  });

  test('showSearchがfalseの時に検索フィールドを非表示にする', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: mockTransactions,
        showSearch: false,
      },
    });

    const searchFields = wrapper.findAllComponents({ name: 'VTextField' });
    const searchField = searchFields.find(field => field.props('label') === '検索');
    expect(searchField).toBeUndefined();
  });

  test('収入金額をプラス記号付きでフォーマットする', () => {
    const incomeTransaction = [
      {
        ...mockTransactions[0],
        type: 'income' as const,
        amount: 50000,
      },
    ];

    const wrapper = mount(TransactionTable, {
      props: {
        transactions: incomeTransaction,
      },
    });

    expect(wrapper.text()).toContain('+¥50,000');
  });

  test('支出金額をマイナス記号付きでフォーマットする', () => {
    const expenseTransaction = [
      {
        ...mockTransactions[1],
        type: 'expense' as const,
        amount: 1200,
      },
    ];

    const wrapper = mount(TransactionTable, {
      props: {
        transactions: expenseTransaction,
      },
    });

    expect(wrapper.text()).toContain('-¥1,200');
  });

  test('取引配列が空の時にデータなしメッセージを表示する', () => {
    const wrapper = mount(TransactionTable, {
      props: {
        transactions: [],
        loading: false,
      },
    });

    expect(wrapper.text()).toContain('取引がありません');
  });
});
