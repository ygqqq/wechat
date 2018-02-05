
<template>
    <mt-header fixed :title="headerTitle">
        <router-link to="/" slot="left">
            <mt-button icon="back">返回</mt-button>
        </router-link>
        <!-- <mt-button icon="more" slot="right" class="hasMoreButton">
        </mt-button> -->
        <mt-button slot="right" class="addFriend" @click="addFriend">
            +
        </mt-button>
    </mt-header>
</template>

<script>
import { Header } from 'mint-ui'
import { MessageBox } from 'mint-ui'
import { getFriends } from '../content/script/getFriends'
import { mapState,mapMutations,mapActions } from "vuex"

const ErrorMsg = 0  //错误消息
const OnlineRemind	= 1	//上线提醒
const OfflineRemind   = 2 //下线提醒 
const AddFriendReq	= 3 //添加好友请求
const AgreeAdd		= 4 //同意好友请求
const DisAgreeAdd 	= 5 //拒绝好友请求
const NormalMsg		= 10 //普通通知消息

export default {
    name:'topHeader',
    data () {
        return {
        }
    },
    props :{
        headerTitle: {
            type: String,
            default: '0',
        }
    },
    methods: {
      ...mapMutations([
        'getCookie', // 将 `this.increment()` 映射为 `this.$store.commit('increment')`
      ]),
      addFriend () {
        MessageBox.prompt('请输入好友姓名').then(({ value, action }) => {
          if( action == 'confirm') {
            if (value == null || value.length <= 0) {
              MessageBox('', '好友昵称不能为空') 
              return
            }
            if (value == this.username) {
              MessageBox('', 'headerTitle不能添加自己') 
              return
            }
            this.$store.state.ws.send(
              JSON.stringify({
                src: this.username,
                dst: value,
                messagetype: AddFriendReq
              }
            ))

            let _this = this
            getFriends(_this,_this.username).then(function(res){
              _this.myFriends = res
              _this.$store.state.userMessage = res
            })
          } 
        },() => {
          console.log(22)
        });
      },
      agreeFriend (otherName) {
        MessageBox.confirm(otherName+'想加你好友，同意吗？').then(() => {
          this.$store.state.ws.send(
            JSON.stringify({
              src: this.username,
              dst: otherName,
              messagetype: AgreeAdd
            }
          ))
        },() => {
          this.$store.state.ws.send(
            JSON.stringify({
              src: this.username,
              dst: otherName,
              messagetype: DisAgreeAdd
            }
          ))
        });
      }
    },
    computed: {
      ...mapState({
        username: 'username',
      })
    },
    created () {
      this.getCookie()
    }
}
</script>

<style lang="scss" scoped type="text/css">
.addFriend{
    font-size: 26px;
}
.mint-header{
    height: 44px;
}
</style>

