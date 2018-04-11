<template>
  <div v-if="authenticating">
    Loading...
  </div>
  <AuthLayout v-else-if="this.$route.meta.group === 'auth'"/>
  <BaseLayout v-else-if="authenticated"/>
</template>

<script>
import {mapState} from 'vuex'

import AuthLayout from './pages/AuthLayout'
import BaseLayout from './pages/BaseLayout'

export default {
  name: 'App',
  components: {
    AuthLayout,
    BaseLayout
  },
  computed: mapState({
    authenticating: (state) => state.auth.authenticating,
    authenticated: (state) => !state.auth.authenticating && !state.auth.authenticateErr
  })
}
</script>

<style>
</style>
