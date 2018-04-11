<template>
  <div :class="itemClasses">
    <div :class="headerClasses" @click="toggle">
      <Icon v-if="!noIcon" :type="icon"/>
      <slot></slot>
      <slot name="extra"></slot>
    </div>
    <collapse-transition>
      <div :class="contentClasses" v-show="isActive">
        <div :class="boxClasses"><slot name="content"></slot></div>
      </div>
    </collapse-transition>
  </div>
</template>
<script>
import Icon from 'iview/src/components/icon/icon'
import CollapseTransition from 'iview/src/components/base/collapse-transition'
const prefixCls = 'ivu-collapse'

export default {
  name: 'MPanel',
  components: { Icon, CollapseTransition },
  props: {
    name: {
      type: String
    },
    disabled: {
      type: Boolean
    },
    noIcon: {
      type: Boolean
    },
    icon: {
      type: String,
      default: 'arrow-right-b'
    }
  },
  data () {
    return {
      index: 0, // use index for default when name is null
      isActive: false
    }
  },
  computed: {
    itemClasses () {
      return [
        `${prefixCls}-item`,
        {
          [`${prefixCls}-item-active`]: this.isActive
        }
      ]
    },
    headerClasses () {
      return `${prefixCls}-header`
    },
    contentClasses () {
      return `${prefixCls}-content`
    },
    boxClasses () {
      return `${prefixCls}-content-box`
    }
  },
  methods: {
    toggle () {
      if (this.disabled) {
        return
      }

      this.$parent.toggle({
        name: this.name || this.index,
        isActive: this.isActive
      })
    }
  }
}
</script>
