<template>
  <div class="auth-wrapper">
    <div class="auth-inner">
      <Logo class="logo" big/>
      <transition appear mode="out-in" name="fade">
        <router-view/>
      </transition>
    </div>
  </div>
</template>

<script>
import Logo from '../components/Logo'

export default {
  name: 'AuthLayout',
  components: {Logo}
}
</script>

<style lang="less">
  @import '../assets/variables.less';

  .auth-wrapper {
    h1 {
      margin-top: 22px;
      margin-bottom: 12px;
      font-size: 22px;
    }

    .invite {
      font-size: 15px;
      font-weight: 500;
      margin: 0 0 15px;
    }

    .actions {
      text-align: left;
    }

    .server-error {
      color: @error-color;
      font-size: 12px;
      float: right;
      line-height: 32px;
    }
  }
</style>

<style scoped>
  .auth-inner {
    margin: 40px auto;
    max-width: 500px;
    text-align: center;
    background: #FFFFFF;
    box-shadow: 0 2px 40px 0 rgba(0, 0, 0, 0.17);
    border-radius: 4px;
    padding: 30px 70px;
    transition: height 0.50s ease 0.2s
  }
</style>
