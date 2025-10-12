import { ref } from 'vue';

interface NotificationOptions {
  message: string;
  type?: 'success' | 'error' | 'warning' | 'info';
  timeout?: number;
}

const show = ref(false);
const message = ref('');
const color = ref('info');
const timeout = ref(3000);

export function useNotification() {
  const showNotification = (options: NotificationOptions) => {
    message.value = options.message;
    timeout.value = options.timeout || 3000;

    // Map type to Vuetify color
    switch (options.type) {
      case 'success':
        color.value = 'success';
        break;
      case 'error':
        color.value = 'error';
        break;
      case 'warning':
        color.value = 'warning';
        break;
      case 'info':
      default:
        color.value = 'info';
        break;
    }

    show.value = true;
  };

  const success = (msg: string) => {
    showNotification({ message: msg, type: 'success' });
  };

  const error = (msg: string) => {
    showNotification({ message: msg, type: 'error' });
  };

  const warning = (msg: string) => {
    showNotification({ message: msg, type: 'warning' });
  };

  const info = (msg: string) => {
    showNotification({ message: msg, type: 'info' });
  };

  return {
    show,
    message,
    color,
    timeout,
    showNotification,
    success,
    error,
    warning,
    info,
  };
}
