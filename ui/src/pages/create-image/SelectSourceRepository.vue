<template>
  <Row class="repository-selector" type="flex">
    <Col class="owners">
      <p><strong>OWNER</strong></p>
      <ul class="owner-list">
        <li class="owner-item"
            v-for="owner in owners.data"
            :key="owner"
            :class="{active: owner === image.repository.owner}"
            @click="selectOwner(owner)">
          {{owner}}
        </li>
      </ul>
      <Spin fix v-if="owners.loading"/>
    </Col>
    <Col class="repositories">
      <div class="header">
        <p class="title"><strong>REPOSITORY</strong></p>
        <div class="search">
          <Input v-model="repositoryFilter" placeholder="Filter repositories" icon="search"/>
        </div>
      </div>
      <div class="repository-list">
        <div class="repository-item" v-for="repo in filteredRepositories" :key="repo" @click="selectRepository(repo)">
          <span class="repository-name">{{repo}}</span>
          <Button class="select-button" size="small" type="text">SELECT</Button>
        </div>
      </div>
      <Spin fix v-if="repositories.loading"/>
    </Col>
  </Row>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'SelectSourceRepository',
  props: {
    next: {type: Function, required: true}
  },
  data () {
    return {
      repositoryFilter: ''
    }
  },
  computed: {
    ...mapState({
      image: (state) => state.newImage.data,
      owners: (state) => state.newImage.owners,
      repositories: (state) => state.newImage.repositories
    }),
    filteredRepositories () {
      const reposByOwner = this.repositories.data[this.image.repository.owner] || []
      if (this.repositoryFilter === '') {
        return reposByOwner
      }

      const nameRegExp = new RegExp(this.repositoryFilter)
      return reposByOwner.filter(repo => nameRegExp.test(repo))
    }
  },
  methods: {
    selectOwner (owner) {
      this.$emit('update:image', {...this.image, repository: {...this.image.repository, owner}})
    },
    selectRepository (repo) {
      this.$emit('update:image', {...this.image, repository: {...this.image.repository, name: repo}})
      this.next()
    }
  }
}
</script>

<style lang="less" scoped>
  @import '../../assets/variables';

  .repository-selector {
    margin: -16px;
    background: #F3F7FF;
  }

  p {
    font-size: @font-size-base - 1;
  }

  .owners {
    background: #fff;
    min-width: 200px;
    height: 300px;
    border-right: 1px solid #D3D9E6;
    max-height: 300px;
    overflow-y: auto;

    p {
      padding: 21px 25px 16px;
    }

    .owner-list {
      list-style: none;
      margin: 0;
      padding: 0;
    }

    .owner-item {
      display: block;
      padding: 13px 25px;
      font-size: @font-size-base;
      cursor: pointer;

      &:hover {
        background: lighten(@primary-color, 40%)
      }

      &.active {
        background: @primary-color;
        color: #fff;
      }
    }
  }

  .repositories {
    padding: 21px 25px;
    flex-grow: 1;
    max-height: 300px;
    overflow-y: auto;

    .header {
      display: flex;

      .search {
        margin-left: auto;
        margin-top: -8px;
        max-width: 250px;
      }
    }

    .repository-list {
      margin-top: 14px;
    }

    .repository-item {
      display: flex;
      padding: 8px 13px;
      border: 1px solid #D3D9E6;
      border-radius: 4px;
      background: #FFFFFF;
      font-size: @font-size-base;
      margin-bottom: 10px;
      cursor: pointer;

      .repository-name {
        flex-grow: 1;
        line-height: 24px;
      }

      .select-button {
        display: none;
        margin-left: auto;
      }

      &:hover {
        border-color: @primary-color;
        color: @primary-color;

        .select-button {
          display: block;
          color: @primary-color;
        }
      }
    }
  }
</style>
