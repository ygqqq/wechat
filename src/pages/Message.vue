<template>
    <div>
      <topHeader :headerTitle="username"></topHeader>
        <div class="message-box">
          <ul v-if="userFriends.length > 0">
            <router-link tag="li" v-for="item in userFriends" :key="item.Id" :to="'/chat/'+item.NickName" >
              <p class="user-img">
                <i>{{getName(item.NickName)}}</i>
              </p>
              <p class="user-mess">
                <span>{{item.NickName}}</span>
                <i>你的最新一条消息</i>
              </p>
              <em class="friend-state" :class="{ 'on' : item.Status? true : false}"></em>
            </router-link>                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            
          </ul>
          <div v-if="userFriends.length <= 0">
            好友列表为空，可点击右上角添加好友
          </div> 
        </div>

    </div>
  </template>
  
  <script>
  import Vue from "vue"
  import axios from "axios"
  import config from "../../config/local.config"
  import topHeader from "../components/header"
  import { Header } from 'mint-ui'
  import { MessageBox } from 'mint-ui'
  import { getFriends } from '../content/script/getFriends'
  import { mapState,mapGetters,mapMutations,mapActions } from "vuex"

  const ErrorMsg = 0  //错误消息
  const OnlineRemind	= 1	//上线提醒
  const OfflineRemind   = 2 //下线提醒 
  const AddFriendReq	= 3 //添加好友请求
  const AgreeAdd		= 4 //同意好友请求
  const DisAgreeAdd 	= 5 //拒绝好友请求
  const ChatMsg			= 6 //普通聊天消息
  const NormalMsg		= 10 //普通通知消息

  export default{
    name: 'message',
    data () {
      return {
        nowTitle: '消息',  
        headerTitle: this.username
      }
    },
    components: {topHeader},
    computed:{
      ...mapState([
        'username','userFriends'
      ]),
      title(){
        return this.$store.state.userFriends
      }
    }, 
    watch:{
      userFriends: function (val){
        this.$store.state.userFriends = val
      }
    },
    methods: {
      getName (userName) {
        return userName.substr(userName.length-1,1)
      }
    },
    mounted () {
      // let storeUserMessage = this.$store.state.userMessage
      // if (!storeUserMessage) {
      //   getFriends(_this,_this.username).then(function(res){
      //     storeUserMessage = res
      //   })
      // }
      

          // ws.addEventListener('message', function(e) {  //监听WebSocket
          //   var msg = JSON.parse(e.data);
          //   console.log(msg)
          //   switch(msg.messagetype){
          //     case ErrorMsg:     
          //       MessageBox('', msg.message)
          //       break
          //     case AddFriendReq:     
          //       _this.$store.commit('agreeFriend', msg.src)
          //       break
          //     case NormalMsg:   
          //       MessageBox('', msg.message)
          //       if (msg.message === '添加成功') {
          //         getFriends(_this,_this.username).then(function(res){
          //           storeUserMessage = res
          //         })
          //       }
          //       break
          //     case OnlineRemind:
          //     case OfflineRemind:
          //       _this.$store.commit('setState', msg)
          //     case NormalMsg:
          //       _this.$store.commit('reception', msg)
          //   }
          // })
    }
  }
  </script>
  
  <style lang="scss" scoped type="text/css">
  
    .addFriend{
      font-size: 26px;
    }
    .message-box{
      padding: 44px 0 0;

      ul {
        
        li{
          padding: 0 10px;
          line-height: 50px;
          height: 50px;
          clear: both;
          display: flex;
          background: #fff;

          .user-img{
            width: 36px;
            height: 36px;
            background: #3296fa;
            border-radius: 50%;
            color: #fff;
            line-height: 36px;
            margin-top: 7px;
            flex:0 0 36px;
          }

          .user-mess{
            flex: 1;
            border-bottom: 1px solid #f3f3f4;
            text-align: left;
            padding: 6px 0 6px 10px;
            display: flex;
            flex-direction: column;

            span{
              display: block;
              font-size: 14px;
              line-height: 14px;
              margin-bottom: 12px;
            }
            i{
              font-size: 10px;
              line-height: 10px;
              color: #abadaf;
            }
          }

          .friend-state{
            width: 8px;
            height: 8px;
            flex:0 0 8px;
            margin-top: 22px;
            background: #bbb;
            border-radius: 50%;

            &.on{
              background: #68bd6e;
            }
          }
        }
      }
    }
  </style>