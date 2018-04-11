<template>
  <div>
    <Table :data="builds.data.builds" :columns="buildColumns"/>
    <Spin fix v-if="builds.loading"/>
  </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'History',
  computed: {
    ...mapState(['builds'])
  },
  data () {
    return {
      buildColumns: [
        {title: 'Image Name', render: (h, obj) => <span>{this.builds.data.images[obj.row.image_key].name}</span>},
        {title: 'Started At', render: (h, obj) => <span>{obj.row.records[0].started_at}</span>},
        {title: 'Finished At', render: (h, obj) => <span>{obj.row.records[0].status === 'in_progress' ? '-' : obj.row.records[0].finished_at}</span>},
        {title: 'Docker Tag', render: (h, obj) => <span>{obj.row.records[0].tag}</span>},
        {title: ' ', align: 'right', render: (h, obj) => <router-link to={{name: 'image-history', params: {key: obj.row.image_key}}}>All Records</router-link>}
      ]
    }
  },
  created () {
    this.$store.dispatch('loadBuilds')
  }
}
</script>

<style scoped>

</style>
