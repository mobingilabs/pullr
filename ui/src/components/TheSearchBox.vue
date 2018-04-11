<template>
  <div class="search-box">
    <AutoComplete placeholder="Type to search" icon="search" @on-search="search" :data="autocompleteData"/>
    <Spin v-if="loadingImages" fix/>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'

export default {
  name: 'TheSearchBox',
  computed: {
    ...mapState({
      loadingImages: (state) => state.images.loading
    }),
    ...mapGetters(['filteredImages'])
  },
  data () {
    return {
      autocompleteData: []
    }
  },
  methods: {
    search (query) {
      this.autocompleteData = this.filteredImages(query).map(img => img.name)
    }
  }
}
</script>

<style lang="less" scoped>
  .search-box {
    .ivu-auto-complete {
      max-width: 250px;
      margin-left: 20%;
    }

    .search-input {
      line-height: 100px;
      border: 0;
    }
  }
</style>
