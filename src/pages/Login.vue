<template>

  <div class="login">
    <div class="login-box">
      <div class="login-box-flex">
        <mt-field label="用户名" placeholder="请输入用户名" v-model="username"></mt-field>
        <mt-field label="密码" placeholder="请输入密码" type="password" v-model="password"></mt-field>
        <div class="oprate-area">
          <mt-button type="primary" size="large" @click="login">登陆</mt-button>
          <div class="register">
            没有账号? <router-link  to="/register">立即注册</router-link>
          </div>
        </div>
      </div>
    </div>
    <div class="login-bgbox">
      <img src="../content/images/bgLogin.jpg"/>
    </div>
    <div class="login-bgboxaa"></div>
  </div>
</template>

<script>
import axios from "axios"
export default {
  data () {
    return {
      username: '',
      password: ''
    }
  },
  methods:{
    login(){
        axios.post('/api/user/login', {
            username: this.username,
            password: this.password
        })
        .then(function (response) {
            let res = response.data;
                alert(res.msg)
        })
        .catch(function (error) {
            console.log(error);
        });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style lang="scss" scoped type="text/css">

.register{
  color: #666;
  font-size: 14px;
  margin-top: 50px;
  display: block;
}
.oprate-area{
  margin-top: 35px;
  text-align: center;
}

.login{
  height: 100%;
  width: 100%;
  position: relative;
  overflow: hidden;

  .login-bgbox{
    filter: blur(55px);
    // background-image: url('../content/images/bgLogin.jpg');
    background-repeat: no-repeat;
    background-size:150% 150%;
    background-attachment: fixed;
    background-position: center;
    position: absolute;
  }
  .login-bgboxaa{
    width: 100%;
    height: 100%;
    display: block;
    position: absolute;
    z-index: 2;
  }
  .login-box{
    position: absolute;
    z-index: 3;
    width: 100%;
    height: 100%;
    display: flex;

    .login-box-flex{
      width: 84%;
      margin: auto;

      .mint-cell{
        background: rgba(-0,0,0, 0.5);
        color: #fff;
        margin: 20px 0;
        border-radius: 5px;
        padding: 4px 0;

        .mint-cell-wrapper{
          background-image: none;
        }
      }
    }
  }
}

</style>
