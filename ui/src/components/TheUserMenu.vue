<template>
  <Dropdown class="user-menu" trigger="click" @on-click="handleUserAction">
    Welcome
    <a href="javascript:void(0)">
      {{user.username}}
      <Icon type="arrow-down-b"></Icon>
    </a>
    <DropdownMenu slot="list">
      <DropdownItem name="settings">Settings</DropdownItem>
      <DropdownItem name="logout">Logout</DropdownItem>
    </DropdownMenu>
  </Dropdown>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'TheUserMenu',
  computed: mapState({
    user: (state) => state.auth.profile.data.user
  }),
  methods: {
    async handleUserAction (action) {
      switch (action) {
        case 'logout':
          await this.$store.dispatch('logout')
          this.$router.push({name: 'login'})
          break
        default:
          break
      }
    }
  }
}
</script>

<style lang="less" scoped>
  @import "../assets/variables.less";

  .user-menu {
    color: lighten(@text-color, 30%);
    a {
      font-weight: 500;
      color: @text-color;
    }
  }
</style>
