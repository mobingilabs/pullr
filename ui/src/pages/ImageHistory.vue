<template>
  <div class="image-history">
    <h1>Build History of {{image.name}}</h1>
    <Table class="records-table" :data="records" :columns="recordColumns"></Table>
    <Page class="pagination" :current="currentPage + 1" :total="pagination.total" :page-size="perPage" @on-page-size-change="onPerPageChanged($event)" @on-change="onPageChanged" show-sizer show-total></Page>
    <Spin fix v-if="loading" />
  </div>
</template>

<script>
import Vue from 'vue'
export default {
  name: 'ImageHistory',
  data () {
    return {
      loading: false,
      records: [],
      currentPage: 0,
      perPage: 20,
      pagination: {total: 0, last: 0, next: 0},
      recordColumns: [
        {title: 'Docker Tag', key: 'tag'},
        {title: 'Started At', render: (h, obj) => <span>{Vue.filter('date')(obj.row.started_at)}</span>},
        {title: 'Finished At', render: (h, obj) => <span>{obj.row.status === 'in_progress' ? '-' : Vue.filter('date')(obj.row.finished_at)}</span>},
        {title: 'Status', key: 'status'}
      ]
    }
  },
  computed: {
    image () {
      return this.$store.getters.getImageByKey(this.$route.params.key) || {}
    }
  },
  methods: {
    onPerPageChanged (perPage) {
      this.perPage = perPage
      this.load(this.currentPage)
    },
    onPageChanged (page) {
      this.load(page)
    },
    async load (page) {
      this.loading = true
      try {
        const data = await this.$store.dispatch('loadImageBuilds', {key: this.$route.params.key, page: this.currentPage, perPage: this.perPage})
        this.records = data.build_records
        this.pagination = data.pagination
        if (page > this.pagination.last) {
          this.currentPage = page - 1
        } else {
          this.currentPage = page
        }
      } catch (err) {
        this.$Notice.error({title: 'Failed to load history', desc: err.message})
      } finally {
        this.loading = false
      }
    }
  },
  created () {
    this.$store.dispatch('loadImage', this.$route.params.key)
    this.load(0)
  }
}
</script>

<style scoped>
  .records-table {
    margin: 12px 0;
  }

  .image-history .pagination {
    text-align: right;
  }
</style>
