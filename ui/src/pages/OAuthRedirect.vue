<template>
  <div>
    <h1 class="title">{{provider}} Login</h1>
    <p>Successfully logged in. Please wait you will be redirected to where you left.</p>
  </div>
</template>

<script>
export default {
  name: 'OAuthRedirect',
  computed: {
    provider () {
      return this.$route.query.provider || ''
    }
  },
  mounted () {
    const opener = window.opener ? window.opener : window['dialogArguments']

    opener.postMessage('OAUTH_SUCCESS', location.origin)
    setTimeout(() => window.close(), 1000)
  }
}
</script>

<style scoped>

</style>
