<!-- 问卷中问题的添加编辑删除-->
<template>
    <div class="Design" v-loading="loading" element-loading-text="加载中，请稍等...">
        <h3>{{title}}</h3>
        <div class="top" v-if="desc != ''">
            {{desc}}
        </div>
        <el-card class="box-card" v-for="(item,index) in QuesAndOp" :key="(index+1).toString()" style="margin 10px;">
            <div slot="header" class="clearfix">
                <div class="quesTitle">
                <span style="color: black; margin-right:3px;"> {{(index+1)+'.'}}  </span>
                {{item.title}}
                <span v-show="item.type==3">&nbsp;&nbsp;&nbsp;&nbsp;(MinValue:&nbsp;{{item.DataMin}}&nbsp;&nbsp;&nbsp;&nbsp;MaxValue:&nbsp;{{item.DataMax}})</span>
                </div>
                <div style="float:right;">
                    <el-button style="padding:2px" type="text" @click="editQues(item)">编辑</el-button>
                    <el-button style="padding:2px" type="text" @click="DelQues(item.QuesId,item.title)">删除</el-button>
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
            <template v-if=" item.type == 3">
                <el-input-number v-model="item.textvalue" style="width:80%" :step="0.01" :precision="2" :min="item.DataMin" :max="item.DataMax"></el-input-number>
            </template>
        </el-card>
        <el-button icon="el-icon-circle-plus" @click="AddQues" style="margin-top:10px;">添加题目</el-button>
            <br>
            <br>
            <br>
        <el-dialog :title="AddTitle" :visible.sync="QuesAddShow" :close-on-click-modal="false" class="AddQues">
            <el-form ref="form" :model="AddQuesModel" label-width="80px">
                <el-form-item label="题目类型" style="width:100%">
                    <el-select v-model="typeValue" placeholder="请选择问题类型" @change="typeSelect">
                        <el-option v-for="(item,index) in allType" :key="item.value" :label="item.label" :value="(index+1)"></el-option>
                    </el-select>
                </el-form-item>
                <el-form-item label="题目标题" style="width:100%">
                    <el-input v-model="AddQuesModel.title" placeholder="请输入题目标题"></el-input>
                </el-form-item>
                <template v-if="AddQuesModel.type==1 || AddQuesModel.type==2">
                    <el-form-item  v-for="(item,index) in AddQuesModel.options" :key="(index+1).toString()" :label="'选项'+(index+1)">
                        <el-row>
                            <el-col :span="16">
                                <el-input v-model="item.title" placeholder="请输入选项" style="width:90%"></el-input>
                            </el-col>
                            <el-col :span="8">
                                <el-button type="danger" plain @click="DelOp(index)">删除选项</el-button>
                            </el-col>
                        </el-row>
                    </el-form-item>
                    <el-button type="primary" plain @click="AddOp">增加选项</el-button>
                </template>
                <template v-if="AddQuesModel.type == 3">
                    <el-form-item label="取值下限">
                        <el-input-number v-model="AddQuesModel.DataMin" style="width:80%" :step="0.01" :precision="2"></el-input-number>
                    </el-form-item>
                    <el-form-item label="取值上限">
                        <el-input-number v-model="AddQuesModel.DataMax" style="width:80%" :step="0.01" :precision="2"></el-input-number>
                    </el-form-item>
                    <el-form-item label="填空">
                        <el-input type="textarea" style="width: 80%" ></el-input>
                    </el-form-item>
                </template>
            </el-form>
            <br>
            <br>
            <div style="width:100%;text-align:right;">
                <el-button style="margin-left:10px;" @click="AddQuesShow=false">取消</el-button>
                <el-button type="primary" style="margin-left:10px;" @click = "AddQuesCheck">确认</el-button>
            </div>
        </el-dialog>


    </div>
</template>

<script>
import * as API from '../api/question'
export default {
    data(){
        return{
            QuesAddShow:false,
            AddTitle:'',
            QuesAndOp:[],
            WjId:0,
            title:'',
            desc:'',
            loading:false,
            typeValue:'',
            AddQuesModel:{
                WjId:0,
                options:[{
                    title:'',
                    opId:0,
                }],
                QuesId:0,
                type:0,
                title:'',
                DataMin:0.0,
                DataMax:0.0,
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
            DelArg:{
                WjId:0,
                QuesId:0,
            }
        }
    },
    methods: {
        init(wjId,title,desc){
            this.WjId=wjId;
            this.title=title;
            this.desc=desc;
            this.getQues(); 
        },
        getQues(){
            this.QuesAndOp=[];
            this.loading=true;
            var Formdata = new FormData();
            Formdata.append('WjId',this.WjId.toString());
            API.ShowQuestions(Formdata).then(res=>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.QuesAndOp = res.data;
                    console.log(this.QuesAndOp);
                    this.loading=false;
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
            })
        },
        AddQues(){
            
            if(this.WjId==0 || this.WjId==null){
                this.$message({
                    type:"error",
                    message:"请先创建问卷",
                })
            }
            this.AddTitle='添加题目';
            this.AddQuesModel={
                WjId:0,
                options:[{
                    title:'',
                    opId:0,
                }],
                QuesId:0,
                type:0,
                title:'',
                radioValue:'',
                checkboxValue:[],
                textvalue:'',
                DataMin:0.0,
                DataMax:0.0,
            };
            this.AddQuesModel.WjId=this.WjId;
            this.QuesAddShow=true;
        },
        AddQuesCheck(){
            //console.log(this.AddQuesModel.WjId);
            console.log(this.AddQuesModel.DataMin)
            console.log(this.AddQuesModel.DataMax)
            this.AddQuesModel.DataMin = parseFloat(this.AddQuesModel.DataMin)
            this.AddQuesModel.DataMax = parseFloat(this.AddQuesModel.DataMax)
            console.log(this.AddQuesModel.DataMin)
            console.log(this.AddQuesModel.DataMax)
            //console.log(this.AddQuesModel)
            API.UpdateQuestion(this.AddQuesModel).then(res=>{
                if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.getQues();
                }else{
                    this.$message({
                        type:"error",
                        message:res.msg+'!',
                    })
                }
                this.QuesAddShow=false;
                this.AddTitle='';
            }).catch(error=>{
                this.$notify.error({
                    message:error,
                })
            })
        },
        editQues(item){
            this.AddQuesModel=item;
            this.AddTitle='编辑问题';
            this.QuesAddShow=true;
        },
        DelQues(QuesId,title){
            this.$confirm('确认删除'+title+'?(删除后将无法恢复)','删除确认',{
                confirmButtonText:'确认删除',
                cancelButtonText:'取消删除',
                type:'warning'
            }).then(()=>{
                this.DelArg.QuesId=QuesId;
                this.DelArg.WjId=0;
                API.DeleteQuestionaire({data:this.DelArg}).then(res=>{
                    if(res.code==200){
                    this.$message({
                        type:"success",
                        message:res.msg+'!',
                    })
                    this.getQues();
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
                })
            })
        },
        AddOp(){
            this.AddQuesModel.options.push({
                title:'',
                opId:0,
            });
        },
        DelOp(index){
            this.AddQuesModel.options.splice(index,1);
        },
        typeSelect(value){
            this.AddQuesModel.type=value;
        },

        changeInputPt2 (e) {
            //console.log(e.target.value)
            if ((e.target.value.indexOf('.') >= 0)) {
              e.target.value = e.target.value.substring(0, e.target.value.indexOf('.') + 3) // 这里采用截取 既可以限制第三位小数输入 也解决了数字输入框默认四舍五入问题
            }
            //console.log(e.target.value)
        },
    },
}
</script>