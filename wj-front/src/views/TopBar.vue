<!-- 顶栏 -->
<template>
  <div class="main"> 
    <el-container>
      <el-header>
        <div style="float:left;margin-left: 50px;line-height:60px;"   @click="toHome">
          <span style="font-size:25px;">在线问卷</span>
        </div>
        <div style="float:right;margin-right: 50px;line-height:60px;">
          <!--未登录状态-->
          <template v-if="!showuser">
            <el-button type="primary" plain style="font-size:15px;" @click="toLogin">登录</el-button>
            <el-button plain style="font-size:15px;" @click="toRegister">注册</el-button>
          </template>
          <!--已登录状态-->
          <template v-else>
            <el-dropdown @command="handleCommand">
              <span>
                {{username}}<i class="el-icon-arrow-down el-icon--right"></i>
              </span>
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item command="a">问卷列表</el-dropdown-item>
                <el-dropdown-item command="b">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
          </template>
        </div>
      </el-header>
      <el-main style="padding:0">
        <router-view @state="state" />
      </el-main>
    </el-container>
  </div>
</template>

<script>
import * as API from "../api/user"
export default {
  name:'TopBar',
  data(){
    return {
      showuser: false,
      username:'',
    }
  },
  methods: {
    toHome(){
      if(showuser==true){
        this.$router.push({path:'/home'});
      }else{
        this.$router.push({path:'/Login'});
      }
    },
    toLogin(){
      this.$router.push({path:'/Login'});
    },
    toRegister(){
      this.$router.push({path:'/Register'});
    },
    state(){
      if(sessionStorage.getItem('username')!=null){
        this.showuser=true;
        this.username=sessionStorage.getItem('username');
        console.log(this.username)
      }else{
        this.showuser=false;
      }
    },
    handleCommand(command){
      if(command=='a'){
        this.toHome();
      }else if(command=='b'){
        this.Logout();
      }
    },
    Logout(){
      API.logout().then(res=>{
        this.$notify({
            type :'sucess',
            message :this.username +res.msg +'!',
        });
        sessionStorage.clear();
        this.state();
        this.toLogin();
      })
    }
  },
}
</script>
<style>
  .main{
    position: absolute;
    width: 100%;
    height: 100%;
  }
  .el-header{
    border-bottom: 2px solid #409EFF;
    background-color: #9befb8;
  }
  .el-dropdown-link {
    cursor: pointer;
    color: #409EFF;
  }
  .el-icon-arrow-down {
    font-size: 12px;
  }
  .demonstration {
    display: block;
    color: #8492a6;
    font-size: 14px;
    margin-bottom: 20px;
  }
</style>