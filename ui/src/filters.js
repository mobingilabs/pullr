import Vue from 'vue'

export function date (value) {
  if (typeof value === 'string') {
    value = new Date(value)
  }

  return value.toLocaleDateString([], {
    month: '2-digit',
    day: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

Vue.filter('date', date)
