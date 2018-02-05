<template>
  <div>
    <topHeader :headerTitle="headerTitle"></topHeader>
    <div class="message-container">
        <div class="message-box others">
            <span>42分钟前</span>
            <div>
                <i>ni</i>
                <p>ssssss</p>
            </div>
        </div>
        <div class="message-box own">
            <span>52分钟前</span>
            <div>
                <i>wo</i>
                <p>ssssss</p>
            </div>
        </div>
    </div>
    <div class="write-container">
        <input type="text" v-model="writeMessage">
        <input type="button" value="发送" @click="sendMessage">
    </div>
  </div>
</template>

<script>
import topHeader from "../components/header"
import { MessageBox } from 'mint-ui'
import { mapState,mapMutations,mapActions } from "vuex"

  const ErrorMsg = 0  //错误消息
  const OnlineRemind	= 1	//上线提醒
  const OfflineRemind   = 2 //下线提醒 
  const AddFriendReq	= 3 //添加好友请求
  const AgreeAdd		= 4 //同意好友请求
  const DisAgreeAdd 	= 5 //拒绝好友请求
  const ChatMsg			= 6 //普通聊天消息
  const NormalMsg		= 10 //普通通知消息

export default {
    name: 'chat',
    components: {topHeader},
    data () {
        return {
            headerTitle: this.$route.params.username,
            writeMessage: '',
            receive: this.$store.state.userMessage.find( item => item.NickName === this.$route.params.username).UserName
        }
    }, 
    methods:{
        ...mapMutations([
            'getCookie', // 将 `this.increment()` 映射为 `this.$store.commit('increment')`
        ]),
        sendMessage () {
            this.$store.state.ws.send(
                JSON.stringify({
                    src: this.username,
                    dst: this.receive,
                    message: this.writeMessage,
                    messagetype: ChatMsg
                }
            ))
        }
    },
    computed: {
      ...mapState({
        username: 'username',
      })
    },
    created () {
        console.log('chat')
    }
}
</script>

<style lang="scss" scoped type="text/css">
    .message-container{
        position: absolute;
        top: 44px;
        bottom: 40px;
        width: 100%;
        overflow: hidden;
        display: flex;
        flex-direction: column-reverse;
        padding: 0 10px;

        .message-box{
            margin: 8px 0;

            div{
                display: flex;
            }

            span{
                display: block;
                padding: 6px;
                font-size: 12px;
                color: #999;
            }
            p{
                background: #fff;
                max-width: 74%;
                border-radius: 3px;
                border: 1px solid #efefef;
                flex: 1;
            }
            i{
                width: 36px;
                height: 36px;
                display: inline-block;
                background: #3296fa;
                border-radius: 50%;
                color: #fff;
                line-height: 36px;
                flex: 0 0 36px;
            }

            

            &.others div{
                flex-direction: row;

                p{
                    margin-left: 8px;
                }
            }

            &.own div{
                flex-direction: row-reverse;

                p{
                    margin-right: 8px;
                }
            }
        }
    }
    .write-container{
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        background: #fff;
        height: 40px;
        display: flex;
        padding: 0 10px;

        input{
            // border: 1px solid #efefef;
            height: 24px;
            border-radius: 3px;
            margin-top: 8px;

            &[type='text']{
                border: 1px solid #efefef;
                flex: 1;
                margin-right: 10px;
            }

            &[type='button']{
                outline: none;
                -webkit-appearance: none;
                background: #999;
                color: #fff;
                flex: 0 0 50px;
                border: none;
            }
        }
    }
</style>


