<template>
    <div>
        <el-form :rules="rules" ref="LoginForm" :model="LoginForm" class="LoginContainer">
            <h3 class="LoginTitle">网站登录</h3>
            <el-form-item prop="username">
                <el-input type="text" v-model="LoginForm.username" aria-placeholder="请输入用户名"></el-input>
            </el-form-item>
            <el-form-item prop="password">
                <el-input type="password" v-model="LoginForm.password" aria-placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item style="text-align: center">
                <el-button type="primary" @click="submitLogin">登录</el-button>
            </el-form-item>
        </el-form>
    </div>    
</template>

<script>
export default {
    name:"Login",
    data(){
        return {
            LoginForm:{
                username:'admin',
                password:'admin'
            },
            rules:{
                username:[{required:true,message:'请输入用户名',trigger:'blur'}],
                password:[{required:true,message:'请输入密码',trigger:'blur'}]
            }
        }
    },
    methods:{
        submitLogin:function(){
            this.$refs.LoginForm.validate((valid)=>{
                if(valid){
                    alert('submit!');
                }else{
                    this.$notify.error({
                        title:'错误',
                        message:'登入失败'
                    });
                    return false;
                }
            })
        }
    }
}
</script>

<style >
    .LoginContainer{
        border-radius: 15px;
        background-clip: padding-box;
        margin: 188px auto;
        width: 350px;
        padding: 15px 35px 15px 35px;
        background: #fff;
        border: 1px solid #eaeaea;
        box-shadow: 0 0 25px #cac6c6;
    }
    .LoginTitle{
        margin: 0px auto 40px auto;
        text-align: center;
    }
</style>