<template>
  <div class="layout">
    <Layout>
      <TheHeader/>
      <Content class="app-content">
        <transition mode="out-in" name="fade">
          <router-view></router-view>
        </transition>
      </Content>
    </Layout>
  </div>
</template>

<script>
import TheHeader from '../components/TheHeader'

export default {
  name: 'BaseLayout',
  components: {
    TheHeader
  }
}
</script>

<style scoped>
</style>
