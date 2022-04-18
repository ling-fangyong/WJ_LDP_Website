<!-- 展示问卷以供填写 -->
<template>
    <div class="display">
        <div class="questionaire">
            <h3>{{title}}</h3>
            <div class="top" v-if="desc!=''">
                {{desc}}
            </div>
        </div>
        <el-card class="box-card" v-for="(item,index) in QuesAndOp" :key="(index+1).toString()">
             <div slot="header" class="clearfix">
                <div class="quesTitle">
                <span style="color: black; margin-right:3px;"> {{(index+1)+'.'}}  </span>
                {{item.title}}
                </div>
             </div>
             <template v-if=" item.type == 1 ">
                <div class="text item" v-for="(option,index) in item.options" :key="(index+1).toString()">
                    <el-radio v-model="item.radioValue" :label="index" style="margin: 5px">{{option.title}}</el-radio>
                </div>
            </template>
            <template v-if=" item.type == 2">
                <el-checkbox-group v-model="item.checkboxValue">
                    <div class="text item" v-for="(option,index) in item.options" :key="(index+1).toString()">
                        <el-checkbox :label="index" style="margin: 5px;">{{option.title}}</el-checkbox>
                    </div>
                </el-checkbox-group>
            </template>
        </el-card>
        <el-button type="primary" style="margin: 5px;" @click="submitQues" :loading="Loading">提交</el-button>
    </div>
</template>



<script>
import * as API from '../api/answer'
export default {
    data(){
        return{
            WjId:0,
            title:'',
            desc:'',
            QuesAndOp:[],
            Loading:false,
            epsilon:1,
            SubmitModel:{
                WjId:0,
                QuesAndOp:[],
            },
            allType:[
                {
                    value:'radio',
                    label:'单选题',
                },
                {
                    value:'checkbox',
                    label:'多选题',
                },
                {
                    value:'text',
                    label:'填空题'
                },
            ],
        }
    },
    mounted() {
        var WjId=this.$route.params.id;
        var Formdata = new FormData();
        // console.log(typeof WjId);
        Formdata.append('WjId',WjId);
        // console.log(Formdata.get('WjId'));
        API.GetQuestionaire(Formdata).then(res=>{
            if(res.code==200){
                this.$message({
                    type:"success",
                    message:res.msg+'!',
                })
                this.title = res.data.Questionaire.title;
                this.desc = res.data.Questionaire.desc;
                this.QuesAndOp = res.data.Ques;
                this.WjId = parseInt(WjId);
                // console.log("test");
                // console.log(this.QuesAndOp);
            }else{
                this.$message({
                    type:"error",
                    message:res.msg+'!',
                })
            }
        }).catch(error=>{
            this.$notify.error({
                message:error,
            })
        });
    },
    methods: {
        submitQues(){
            // console.log(this.QuesAndOp);
            // 前端变换
            this.Loading = true;
            for(var i=0,len=this.QuesAndOp.length;i<len;i++){
                if(this.QuesAndOp[i].type == 1){//单选题变换
                    this.QuesAndOp[i].radioValue = this.Opchange(this.QuesAndOp[i].options.length,this.QuesAndOp[i].radioValue);
                }else if(this.QuesAndOp[i].type == 2){
                    var t=[];
                    // console.log("change");
                    // console.log(this.QuesAndOp[i].checkboxValue);
                    for(const item in this.QuesAndOp[i].checkboxValue){
                        
                        t.push(this.Opchange(this.QuesAndOp[i].options.length,item));
                    }
                    console.log(t);
                    this.QuesAndOp[i].checkboxValue = t;
                }
            }
            // console.log("change");
            // console.log(this.QuesAndOp);

            //后端交互
            this.SubmitModel.QuesAndOp = this.QuesAndOp;
            this.SubmitModel.WjId = this.WjId;

            API.SubmitQues(this.SubmitModel).then(res=>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.Loading = false;
                    this.$router.push({path:'/Home'});
                }else{
                    this.Loading = false;
                    this.$message({
                        type:"error",
                        message:res.msg+'!',
                    })
                }
            }).catch(error=>{
                this.$notify.error({
                    message:error,
                })
            });

        },
        Opchange(opnum,value){
            // console.log(opnum);
            // console.log(value);
            var p = Math.exp(this.epsilon)/(Math.exp(this.epsilon)+opnum-1);
            var q = 1/(Math.exp(this.epsilon)+opnum-1);
            // console.log(p);
            // console.log(q);
            if(Math.random()>p-q) value = Math.floor(Math.random()*opnum);
            // console.log("value");
            // console.log(value);
            return parseInt(value);
        }
    },
}
</script>