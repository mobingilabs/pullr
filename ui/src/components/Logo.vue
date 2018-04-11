<template>
  <img v-if="big" :src="logoBig" :srcset="logoBig2x">
  <img v-else :src="logo" :srcset="logo2x">
</template>

<script>
import Logo from '@/assets/logo.png'
import Logo2x from '@/assets/logo@2x.png'
import LogoBig from '@/assets/logo-big.png'
import LogoBig2x from '@/assets/logo-big@2x.png'

export default {
  name: 'TheLogo',
  props: {
    big: {
      type: Boolean
    }
  },
  data () {
    return {
      logoSrc: Logo,
      logo2xSrc: Logo2x,
      logoBigSrc: LogoBig,
      logoBig2xSrc: LogoBig2x
    }
  },

  computed: {
    logo () {
      return this.logoSrc
    },
    logo2x () {
      return `${this.logo2xSrc} 2x`
    },
    logoBig () {
      return this.logoBigSrc
    },
    logoBig2x () {
      return `${this.logoBig2xSrc} 2x`
    }
  }

}
</script>

<style scoped>
  img {
    vertical-align: middle;
  }
</style>
