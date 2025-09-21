import 'vuetify/styles'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import '@mdi/font/css/materialdesignicons.css'

export default createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        colors: {
          primary: '#1976d2',
          secondary: '#424242',
          accent: '#82b1ff',
          error: '#ff5252',
          info: '#2196f3',
          success: '#4caf50',
          warning: '#fb8c00',
        },
      },
      dark: {
        colors: {
          primary: '#2196f3',
          secondary: '#424242',
          accent: '#ff4081',
          error: '#ff5252',
          info: '#2196f3',
          success: '#4caf50',
          warning: '#fb8c00',
        },
      },
    },
  },
  icons: {
    defaultSet: 'mdi',
  },
})