<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useCategoryStore } from '@/stores/category';
import CategoryDialog from '@/components/category/CategoryDialog.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import type { Category, CreateCategoryRequest } from '@/types';
import { useNotification } from '@/composables/useNotification';

const categoryStore = useCategoryStore();
const { loading, error } = categoryStore;
const notification = useNotification();

const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const editingCategory = ref<Category | null>(null);

// 削除確認ダイアログ
const showDeleteConfirm = ref(false);
const deletingCategoryId = ref<number | null>(null);

const incomeCategories = computed(() => categoryStore.getIncomeCategories());
const expenseCategories = computed(() => categoryStore.getExpenseCategories());

const editCategory = (category: Category) => {
  editingCategory.value = category;
  showEditDialog.value = true;
};

const deleteCategory = async (id: number) => {
  deletingCategoryId.value = id;
  showDeleteConfirm.value = true;
};

const confirmDelete = async () => {
  if (deletingCategoryId.value === null) return;

  try {
    await categoryStore.deleteCategory(deletingCategoryId.value);
    notification.success('カテゴリを削除しました');
  } catch (err) {
    notification.error('カテゴリの削除に失敗しました。関連する取引が存在する可能性があります。');
  } finally {
    deletingCategoryId.value = null;
  }
};

const closeDialog = () => {
  showCreateDialog.value = false;
  showEditDialog.value = false;
  editingCategory.value = null;
};

const handleSave = async (data: CreateCategoryRequest) => {
  try {
    if (editingCategory.value) {
      await categoryStore.updateCategory(editingCategory.value.id, data);
      notification.success('カテゴリを更新しました');
    } else {
      await categoryStore.createCategory(data);
      notification.success('カテゴリを作成しました');
    }
    closeDialog();
  } catch (err) {
    notification.error('カテゴリの保存に失敗しました');
  }
};

onMounted(() => {
  categoryStore.fetchCategories();
});
</script>

<template>
  <div class="categories">
    <v-row>
      <v-col cols="12">
        <div class="d-flex justify-space-between align-center mb-6">
          <h2 class="text-h4 font-weight-bold">カテゴリ管理</h2>
          <v-btn
            color="primary"
            prepend-icon="mdi-plus"
            size="large"
            @click="showCreateDialog = true"
          >
            新規カテゴリ
          </v-btn>
        </div>
      </v-col>
    </v-row>

    <v-row>
      <!-- 収入カテゴリ -->
      <v-col cols="12" md="6">
        <v-card elevation="2">
          <v-card-title class="bg-success">
            <v-icon start>mdi-cash-plus</v-icon>
            収入カテゴリ
          </v-card-title>

          <v-card-text>
            <v-progress-linear
              v-if="loading"
              indeterminate
              color="success"
            ></v-progress-linear>

            <v-alert
              v-else-if="error"
              type="error"
              variant="tonal"
            >
              {{ error }}
            </v-alert>

            <v-alert
              v-else-if="incomeCategories.length === 0"
              type="info"
              variant="tonal"
            >
              収入カテゴリがありません
            </v-alert>

            <v-list v-else lines="one">
              <v-list-item
                v-for="category in incomeCategories"
                :key="category.id"
                class="mb-2"
                border
                rounded
              >
                <template #prepend>
                  <v-avatar :color="category.color" size="32"></v-avatar>
                </template>

                <v-list-item-title class="font-weight-medium">
                  {{ category.name }}
                </v-list-item-title>

                <template #append>
                  <div class="d-flex gap-2">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      variant="text"
                      color="primary"
                      @click="editCategory(category)"
                    ></v-btn>
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      @click="deleteCategory(category.id)"
                    ></v-btn>
                  </div>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>

      <!-- 支出カテゴリ -->
      <v-col cols="12" md="6">
        <v-card elevation="2">
          <v-card-title class="bg-error">
            <v-icon start>mdi-cash-minus</v-icon>
            支出カテゴリ
          </v-card-title>

          <v-card-text>
            <v-progress-linear
              v-if="loading"
              indeterminate
              color="error"
            ></v-progress-linear>

            <v-alert
              v-else-if="error"
              type="error"
              variant="tonal"
            >
              {{ error }}
            </v-alert>

            <v-alert
              v-else-if="expenseCategories.length === 0"
              type="info"
              variant="tonal"
            >
              支出カテゴリがありません
            </v-alert>

            <v-list v-else lines="one">
              <v-list-item
                v-for="category in expenseCategories"
                :key="category.id"
                class="mb-2"
                border
                rounded
              >
                <template #prepend>
                  <v-avatar :color="category.color" size="32"></v-avatar>
                </template>

                <v-list-item-title class="font-weight-medium">
                  {{ category.name }}
                </v-list-item-title>

                <template #append>
                  <div class="d-flex gap-2">
                    <v-btn
                      icon="mdi-pencil"
                      size="small"
                      variant="text"
                      color="primary"
                      @click="editCategory(category)"
                    ></v-btn>
                    <v-btn
                      icon="mdi-delete"
                      size="small"
                      variant="text"
                      color="error"
                      @click="deleteCategory(category.id)"
                    ></v-btn>
                  </div>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <CategoryDialog
      v-if="showCreateDialog || showEditDialog"
      :show="showCreateDialog || showEditDialog"
      :category="editingCategory"
      @close="closeDialog"
      @save="handleSave"
    />

    <!-- 削除確認ダイアログ -->
    <ConfirmDialog
      v-model="showDeleteConfirm"
      title="カテゴリの削除"
      message="このカテゴリを削除しますか？関連する取引がある場合は削除できません。"
      confirm-text="削除"
      confirm-color="error"
      @confirm="confirmDelete"
    />
  </div>
</template>

<style scoped>
/* Vuetifyのユーティリティクラスとコンポーネントスタイルを使用 */
.gap-2 {
  gap: 8px;
}
</style>
