<template>
  <div class="wrapper-header">
    <Menu mode="horizontal" :active-name="$route.meta.group" @on-select="gotoRoute">
      <div class="wrapper-header-nav">
        <div class="logo">
          <Logo/>
        </div>
        <div class="nav-menu">
          <MenuItem name="images">
            <Icon type="cube" size="24"/>
            Images
          </MenuItem>
          <MenuItem name="history">
            <Icon type="clock" size="24"/>
            Build History
          </MenuItem>
        </div>
        <div class="search">
          <TheSearchBox />
        </div>
        <div class="user-menu">
          <Button type="primary" icon="plus" @click="gotoRoute('create-image')">ADD IMAGE</Button>
          <div class="v-divider"></div>
          <TheUserMenu />
        </div>
      </div>
    </Menu>
  </div>
</template>

<script>
import {mapState} from 'vuex'

import Logo from './Logo'
import TheSearchBox from './TheSearchBox'
import TheUserMenu from './TheUserMenu'

export default {
  name: 'TheHeader',
  components: {
    TheUserMenu,
    TheSearchBox,
    Logo
  },
  computed: mapState({
    user: (state) => state.auth.profile.data.user
  }),
  methods: {
    gotoRoute (name) {
      this.$router.push({name})
    }
  }
}
</script>

<style lang="less" scoped>
  @import "../assets/variables.less";

  .wrapper-header {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 1000;
    height: 77px;
    line-height: 77px;
  }

  .wrapper-header-nav {
    padding: @layout-header-padding;
    display: flex;
    height: inherit;
  }

  .logo {
    margin-right: 50px;
  }

  .ivu-menu {
    height: inherit;
    line-height: inherit;
  }

  .nav-menu {
    height: inherit;
    /*flex-grow: 1;*/

    .ivu-menu-item .ivu-icon {
      vertical-align: middle;
    }
  }

  .search {
    align-self: center;
    flex-grow: 1;
  }
</style>
