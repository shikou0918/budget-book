<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useTransactionStore } from '@/stores/transaction';
import TransactionDialog from '@/components/transaction/TransactionDialog.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import type { Transaction, CreateTransactionRequest } from '@/types';
import { storeToRefs } from 'pinia';
import TransactionTable from '@/components/transaction/TransactionTable.vue';
import { useNotification } from '@/composables/useNotification';

const transactionStore = useTransactionStore();
const { transactions, loading, error } = storeToRefs(transactionStore);
const notification = useNotification();

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingTransaction = ref<Transaction | null>(null);

// 削除確認ダイアログ
const showDeleteConfirm = ref(false);
const deletingTransactionId = ref<number | null>(null);

const editTransaction = (transaction: Transaction) => {
  editingTransaction.value = transaction;
  showEditDialog.value = true;
};

const deleteTransaction = async (id: number) => {
  deletingTransactionId.value = id;
  showDeleteConfirm.value = true;
};

const confirmDelete = async () => {
  if (deletingTransactionId.value === null) return;

  try {
    await transactionStore.deleteTransaction(deletingTransactionId.value);
    notification.success('取引を削除しました');
  } catch (err) {
    notification.error('取引の削除に失敗しました');
  } finally {
    deletingTransactionId.value = null;
  }
};

const closeDialog = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingTransaction.value = null;
};

const handleSave = async (data: CreateTransactionRequest) => {
  try {
    if (editingTransaction.value) {
      await transactionStore.updateTransaction(editingTransaction.value.id, data);
      notification.success('取引を更新しました');
    } else {
      await transactionStore.createTransaction(data);
      notification.success('取引を作成しました');
    }
    closeDialog();
  } catch (err) {
    notification.error('取引の保存に失敗しました');
  }
};

onMounted(() => {
  transactionStore.fetchTransactions();
});
</script>

<template>
  <div class="transactions">
    <v-row>
      <v-col cols="12">
        <div class="d-flex justify-space-between align-center mb-6">
          <h2 class="text-h4 font-weight-bold">取引管理</h2>
        </div>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card elevation="2">
          <v-card-title class="d-flex justify-space-between align-center">
            <span class="text-h6">取引一覧</span>
            <v-btn color="primary" prepend-icon="mdi-plus" @click="showCreateDialog = true">
              新規取引
            </v-btn>
          </v-card-title>
          <v-card-text>
            <TransactionTable
              :transactions="transactions"
              :loading="loading"
              :error="error"
              :items-per-page="10"
              @edit="editTransaction"
              @delete="deleteTransaction"
            />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <TransactionDialog
      v-if="showCreateDialog || showEditDialog"
      :show="showCreateDialog || showEditDialog"
      :transaction="editingTransaction"
      @close="closeDialog"
      @save="handleSave"
    />

    <!-- 削除確認ダイアログ -->
    <ConfirmDialog
      v-model="showDeleteConfirm"
      title="取引の削除"
      message="この取引を削除しますか？"
      confirm-text="削除"
      confirm-color="error"
      @confirm="confirmDelete"
    />
  </div>
</template>
