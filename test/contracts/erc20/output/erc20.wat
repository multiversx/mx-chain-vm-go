(module
  (type (;0;) (func (param i32 i32) (result i32)))
  (type (;1;) (func (param i32 i32 i32 i32)))
  (type (;2;) (func (result i32)))
  (type (;3;) (func (param i32 i32)))
  (type (;4;) (func (param i32)))
  (type (;5;) (func (param i64) (result i32)))
  (type (;6;) (func (param i32 i32 i32) (result i32)))
  (type (;7;) (func (param i32 i32 i32)))
  (type (;8;) (func (param i64)))
  (type (;9;) (func))
  (import "env" "bigIntGetUnsignedBytes" (func (;0;) (type 0)))
  (import "env" "writeLog" (func (;1;) (type 1)))
  (import "env" "getNumArguments" (func (;2;) (type 2)))
  (import "env" "signalError" (func (;3;) (type 3)))
  (import "env" "getCaller" (func (;4;) (type 4)))
  (import "env" "bigIntNew" (func (;5;) (type 5)))
  (import "env" "bigIntGetUnsignedArgument" (func (;6;) (type 3)))
  (import "env" "bigIntStorageStoreUnsigned" (func (;7;) (type 6)))
  (import "env" "bigIntStorageLoadUnsigned" (func (;8;) (type 6)))
  (import "env" "bigIntFinishUnsigned" (func (;9;) (type 4)))
  (import "env" "getArgument" (func (;10;) (type 0)))
  (import "env" "bigIntCmp" (func (;11;) (type 0)))
  (import "env" "bigIntSub" (func (;12;) (type 7)))
  (import "env" "bigIntAdd" (func (;13;) (type 7)))
  (import "env" "int64finish" (func (;14;) (type 8)))
  (func (;15;) (type 1) (param i32 i32 i32 i32)
    (local i32 i32 i32 i32 i32 i32 i32 i32)
    i32.const 0
    local.get 0
    i32.load8_u
    i32.store8 offset=1152
    local.get 0
    i32.load8_u offset=1
    local.set 4
    local.get 0
    i32.load8_u offset=2
    local.set 5
    local.get 0
    i32.load8_u offset=3
    local.set 6
    local.get 0
    i32.load8_u offset=4
    local.set 7
    local.get 0
    i32.load8_u offset=5
    local.set 8
    local.get 0
    i32.load8_u offset=6
    local.set 9
    local.get 0
    i32.load8_u offset=7
    local.set 10
    local.get 0
    i32.load8_u offset=8
    local.set 11
    i32.const 0
    local.get 0
    i32.load8_u offset=9
    i32.store8 offset=1161
    i32.const 0
    local.get 11
    i32.store8 offset=1160
    i32.const 0
    local.get 10
    i32.store8 offset=1159
    i32.const 0
    local.get 9
    i32.store8 offset=1158
    i32.const 0
    local.get 8
    i32.store8 offset=1157
    i32.const 0
    local.get 7
    i32.store8 offset=1156
    i32.const 0
    local.get 6
    i32.store8 offset=1155
    i32.const 0
    local.get 5
    i32.store8 offset=1154
    i32.const 0
    local.get 4
    i32.store8 offset=1153
    i32.const 0
    local.get 0
    i32.load8_u offset=10
    i32.store8 offset=1162
    i32.const 0
    local.get 0
    i32.load8_u offset=11
    i32.store8 offset=1163
    i32.const 0
    local.get 0
    i32.load8_u offset=12
    i32.store8 offset=1164
    i32.const 0
    local.get 0
    i32.load8_u offset=13
    i32.store8 offset=1165
    i32.const 0
    local.get 0
    i32.load8_u offset=14
    i32.store8 offset=1166
    i32.const 0
    local.get 0
    i32.load8_u offset=15
    i32.store8 offset=1167
    i32.const 0
    local.get 0
    i32.load8_u offset=16
    i32.store8 offset=1168
    i32.const 0
    local.get 0
    i32.load8_u offset=17
    i32.store8 offset=1169
    i32.const 0
    local.get 0
    i32.load8_u offset=18
    i32.store8 offset=1170
    i32.const 0
    local.get 0
    i32.load8_u offset=19
    i32.store8 offset=1171
    i32.const 0
    local.get 0
    i32.load8_u offset=20
    i32.store8 offset=1172
    i32.const 0
    local.get 0
    i32.load8_u offset=21
    i32.store8 offset=1173
    i32.const 0
    local.get 0
    i32.load8_u offset=22
    i32.store8 offset=1174
    i32.const 0
    local.get 0
    i32.load8_u offset=23
    i32.store8 offset=1175
    i32.const 0
    local.get 0
    i32.load8_u offset=24
    i32.store8 offset=1176
    i32.const 0
    local.get 0
    i32.load8_u offset=25
    i32.store8 offset=1177
    i32.const 0
    local.get 0
    i32.load8_u offset=26
    i32.store8 offset=1178
    i32.const 0
    local.get 0
    i32.load8_u offset=27
    i32.store8 offset=1179
    i32.const 0
    local.get 0
    i32.load8_u offset=28
    i32.store8 offset=1180
    i32.const 0
    local.get 0
    i32.load8_u offset=29
    i32.store8 offset=1181
    i32.const 0
    local.get 0
    i32.load8_u offset=30
    i32.store8 offset=1182
    i32.const 0
    local.get 0
    i32.load8_u offset=31
    i32.store8 offset=1183
    i32.const 0
    local.get 1
    i32.load8_u
    i32.store8 offset=1184
    i32.const 0
    local.get 1
    i32.load8_u offset=1
    i32.store8 offset=1185
    i32.const 0
    local.get 1
    i32.load8_u offset=2
    i32.store8 offset=1186
    i32.const 0
    local.get 1
    i32.load8_u offset=3
    i32.store8 offset=1187
    i32.const 0
    local.get 1
    i32.load8_u offset=4
    i32.store8 offset=1188
    i32.const 0
    local.get 1
    i32.load8_u offset=5
    i32.store8 offset=1189
    i32.const 0
    local.get 1
    i32.load8_u offset=6
    i32.store8 offset=1190
    i32.const 0
    local.get 1
    i32.load8_u offset=7
    i32.store8 offset=1191
    i32.const 0
    local.get 1
    i32.load8_u offset=8
    i32.store8 offset=1192
    i32.const 0
    local.get 1
    i32.load8_u offset=9
    i32.store8 offset=1193
    i32.const 0
    local.get 1
    i32.load8_u offset=10
    i32.store8 offset=1194
    i32.const 0
    local.get 1
    i32.load8_u offset=11
    i32.store8 offset=1195
    i32.const 0
    local.get 1
    i32.load8_u offset=12
    i32.store8 offset=1196
    i32.const 0
    local.get 1
    i32.load8_u offset=13
    i32.store8 offset=1197
    i32.const 0
    local.get 1
    i32.load8_u offset=14
    i32.store8 offset=1198
    i32.const 0
    local.get 1
    i32.load8_u offset=15
    i32.store8 offset=1199
    i32.const 0
    local.get 1
    i32.load8_u offset=16
    i32.store8 offset=1200
    i32.const 0
    local.get 1
    i32.load8_u offset=17
    i32.store8 offset=1201
    i32.const 0
    local.get 1
    i32.load8_u offset=18
    i32.store8 offset=1202
    i32.const 0
    local.get 1
    i32.load8_u offset=19
    i32.store8 offset=1203
    i32.const 0
    local.get 1
    i32.load8_u offset=20
    i32.store8 offset=1204
    i32.const 0
    local.get 1
    i32.load8_u offset=21
    i32.store8 offset=1205
    i32.const 0
    local.get 1
    i32.load8_u offset=22
    i32.store8 offset=1206
    i32.const 0
    local.get 1
    i32.load8_u offset=23
    i32.store8 offset=1207
    i32.const 0
    local.get 1
    i32.load8_u offset=24
    i32.store8 offset=1208
    i32.const 0
    local.get 1
    i32.load8_u offset=25
    i32.store8 offset=1209
    i32.const 0
    local.get 1
    i32.load8_u offset=26
    i32.store8 offset=1210
    i32.const 0
    local.get 1
    i32.load8_u offset=27
    i32.store8 offset=1211
    i32.const 0
    local.get 1
    i32.load8_u offset=28
    i32.store8 offset=1212
    i32.const 0
    local.get 1
    i32.load8_u offset=29
    i32.store8 offset=1213
    i32.const 0
    local.get 1
    i32.load8_u offset=30
    i32.store8 offset=1214
    i32.const 0
    local.get 1
    i32.load8_u offset=31
    i32.store8 offset=1215
    i32.const 0
    local.get 2
    i32.load8_u
    i32.store8 offset=1216
    i32.const 0
    local.get 2
    i32.load8_u offset=1
    i32.store8 offset=1217
    i32.const 0
    local.get 2
    i32.load8_u offset=2
    i32.store8 offset=1218
    i32.const 0
    local.get 2
    i32.load8_u offset=3
    i32.store8 offset=1219
    i32.const 0
    local.get 2
    i32.load8_u offset=4
    i32.store8 offset=1220
    i32.const 0
    local.get 2
    i32.load8_u offset=5
    i32.store8 offset=1221
    i32.const 0
    local.get 2
    i32.load8_u offset=6
    i32.store8 offset=1222
    i32.const 0
    local.get 2
    i32.load8_u offset=7
    i32.store8 offset=1223
    i32.const 0
    local.get 2
    i32.load8_u offset=8
    i32.store8 offset=1224
    i32.const 0
    local.get 2
    i32.load8_u offset=9
    i32.store8 offset=1225
    i32.const 0
    local.get 2
    i32.load8_u offset=10
    i32.store8 offset=1226
    i32.const 0
    local.get 2
    i32.load8_u offset=11
    i32.store8 offset=1227
    i32.const 0
    local.get 2
    i32.load8_u offset=12
    i32.store8 offset=1228
    i32.const 0
    local.get 2
    i32.load8_u offset=13
    i32.store8 offset=1229
    i32.const 0
    local.get 2
    i32.load8_u offset=14
    i32.store8 offset=1230
    i32.const 0
    local.get 2
    i32.load8_u offset=15
    i32.store8 offset=1231
    i32.const 0
    local.get 2
    i32.load8_u offset=16
    i32.store8 offset=1232
    i32.const 0
    local.get 2
    i32.load8_u offset=17
    i32.store8 offset=1233
    i32.const 0
    local.get 2
    i32.load8_u offset=18
    i32.store8 offset=1234
    i32.const 0
    local.get 2
    i32.load8_u offset=19
    i32.store8 offset=1235
    i32.const 0
    local.get 2
    i32.load8_u offset=20
    i32.store8 offset=1236
    i32.const 0
    local.get 2
    i32.load8_u offset=21
    i32.store8 offset=1237
    i32.const 0
    local.get 2
    i32.load8_u offset=22
    i32.store8 offset=1238
    i32.const 0
    local.get 2
    i32.load8_u offset=23
    i32.store8 offset=1239
    i32.const 0
    local.get 2
    i32.load8_u offset=24
    i32.store8 offset=1240
    i32.const 0
    local.get 2
    i32.load8_u offset=25
    i32.store8 offset=1241
    i32.const 0
    local.get 2
    i32.load8_u offset=26
    i32.store8 offset=1242
    i32.const 0
    local.get 2
    i32.load8_u offset=27
    i32.store8 offset=1243
    i32.const 0
    local.get 2
    i32.load8_u offset=28
    i32.store8 offset=1244
    i32.const 0
    local.get 2
    i32.load8_u offset=29
    i32.store8 offset=1245
    i32.const 0
    local.get 2
    i32.load8_u offset=30
    i32.store8 offset=1246
    i32.const 0
    local.get 2
    i32.load8_u offset=31
    i32.store8 offset=1247
    i32.const 1248
    local.get 3
    i32.const 1248
    call 0
    i32.const 1152
    i32.const 3
    call 1)
  (func (;16;) (type 9)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 1
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 1024
      call 4
      i32.const 0
      i64.const 0
      call 5
      local.tee 1
      call 6
      i32.const 0
      i64.const 0
      i64.store offset=1144
      i32.const 0
      i64.const 0
      i64.store offset=1136
      i32.const 0
      i64.const 0
      i64.store offset=1128
      i32.const 0
      i64.const 0
      i64.store offset=1120
      i32.const 1120
      i32.const 32
      local.get 1
      call 7
      drop
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1024
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1028 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1036 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1044 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1052
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      local.get 1
      call 7
      drop
    end
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;17;) (type 9)
    call 16)
  (func (;18;) (type 9)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.eqz
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i64.const 0
      i64.store offset=1144
      i32.const 0
      i64.const 0
      i64.store offset=1136
      i32.const 0
      i64.const 0
      i64.store offset=1128
      i32.const 0
      i64.const 0
      i64.store offset=1120
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 1
      call 8
      drop
      local.get 1
      call 9
    end
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;19;) (type 9)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 1
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 1088
      call 10
      drop
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1088
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1092 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1100 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1108 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1116
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 1
      call 8
      drop
      local.get 1
      call 9
    end
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;20;) (type 9)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 2
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 1024
      call 10
      drop
      i32.const 1
      i32.const 1056
      call 10
      drop
      i32.const 0
      i32.const 2
      i32.store8 offset=1120
      i32.const 0
      i32.const 0
      i32.load8_u offset=1034
      i32.store8 offset=1121
      i32.const 0
      i32.const 0
      i32.load offset=1035 align=1
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1039 align=1
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1047 align=1
      i32.store16 offset=1134
      i32.const 0
      i32.const 0
      i64.load offset=1066 align=2
      i64.store offset=1136
      i32.const 0
      i32.const 0
      i64.load offset=1074 align=2
      i64.store offset=1144
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 1
      call 8
      drop
      local.get 1
      call 9
    end
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;21;) (type 9)
    (local i32 i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 2
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 1024
      call 4
      i32.const 0
      i32.const 1056
      call 10
      drop
      i32.const 1
      i64.const 0
      call 5
      local.tee 1
      call 6
      block  ;; label = @2
        local.get 1
        i64.const 0
        call 5
        call 11
        i32.const -1
        i32.gt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1368
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1360
        i64.store
        local.get 0
        i32.const 15
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1024
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1028 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1036 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1044 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1052
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 2
      call 8
      drop
      block  ;; label = @2
        local.get 1
        local.get 2
        call 11
        i32.const 1
        i32.lt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i32.load offset=1391 align=1
        i32.store offset=15 align=1
        local.get 0
        i32.const 0
        i64.load offset=1384
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1376
        i64.store
        local.get 0
        i32.const 18
        call 3
        br 1 (;@1;)
      end
      local.get 2
      local.get 2
      local.get 1
      call 12
      i32.const 1120
      i32.const 32
      local.get 2
      call 7
      drop
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1056
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1060 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1068 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1076 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1084
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 2
      call 8
      drop
      local.get 2
      local.get 2
      local.get 1
      call 13
      i32.const 1120
      i32.const 32
      local.get 2
      call 7
      drop
      i32.const 1312
      i32.const 1024
      i32.const 1056
      local.get 1
      call 15
      i64.const 1
      call 14
    end
    local.get 0
    i32.const 32
    i32.add
    global.set 0)
  (func (;22;) (type 9)
    (local i32 i32)
    global.get 0
    i32.const 16
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 2
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 1024
      call 4
      i32.const 0
      i32.const 1056
      call 10
      drop
      i32.const 1
      i64.const 0
      call 5
      local.tee 1
      call 6
      block  ;; label = @2
        local.get 1
        i64.const 0
        call 5
        call 11
        i32.const -1
        i32.gt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1368
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1360
        i64.store
        local.get 0
        i32.const 15
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 2
      i32.store8 offset=1120
      i32.const 0
      i32.const 0
      i32.load8_u offset=1034
      i32.store8 offset=1121
      i32.const 0
      i32.const 0
      i32.load offset=1035 align=1
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1039 align=1
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1047 align=1
      i32.store16 offset=1134
      i32.const 0
      i32.const 0
      i64.load offset=1066 align=2
      i64.store offset=1136
      i32.const 0
      i32.const 0
      i32.load offset=1074 align=2
      i32.store offset=1144
      i32.const 0
      i32.const 0
      i32.load16_u offset=1078
      i32.store16 offset=1148
      i32.const 0
      i32.const 0
      i32.load8_u offset=1080
      i32.store8 offset=1150
      i32.const 0
      i32.const 0
      i32.load8_u offset=1081
      i32.store8 offset=1151
      i32.const 1120
      i32.const 32
      local.get 1
      call 7
      drop
      i32.const 1280
      i32.const 1024
      i32.const 1056
      local.get 1
      call 15
      i64.const 1
      call 14
    end
    local.get 0
    i32.const 16
    i32.add
    global.set 0)
  (func (;23;) (type 9)
    (local i32 i32 i32)
    global.get 0
    i32.const 32
    i32.sub
    local.tee 0
    global.set 0
    block  ;; label = @1
      block  ;; label = @2
        call 2
        i32.const 3
        i32.eq
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1351 align=1
        i64.store offset=7 align=1
        local.get 0
        i32.const 0
        i64.load offset=1344 align=1
        i64.store
        local.get 0
        i32.const 14
        call 3
        br 1 (;@1;)
      end
      i32.const 1088
      call 4
      i32.const 0
      i32.const 1024
      call 10
      drop
      i32.const 1
      i32.const 1056
      call 10
      drop
      i32.const 2
      i64.const 0
      call 5
      local.tee 1
      call 6
      block  ;; label = @2
        local.get 1
        i64.const 0
        call 5
        call 11
        i32.const -1
        i32.gt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i64.load offset=1368
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1360
        i64.store
        local.get 0
        i32.const 15
        call 3
        br 1 (;@1;)
      end
      i32.const 0
      i32.const 2
      i32.store8 offset=1120
      i32.const 0
      i32.const 0
      i32.load8_u offset=1034
      i32.store8 offset=1121
      i32.const 0
      i32.const 0
      i32.load offset=1035 align=1
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1039 align=1
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1047 align=1
      i32.store16 offset=1134
      i32.const 0
      i32.const 0
      i64.load offset=1098 align=2
      i64.store offset=1136
      i32.const 0
      i32.const 0
      i32.load offset=1106 align=2
      i32.store offset=1144
      i32.const 0
      i32.const 0
      i32.load16_u offset=1110
      i32.store16 offset=1148
      i32.const 0
      i32.const 0
      i32.load8_u offset=1112
      i32.store8 offset=1150
      i32.const 0
      i32.const 0
      i32.load8_u offset=1113
      i32.store8 offset=1151
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 2
      call 8
      drop
      block  ;; label = @2
        local.get 1
        local.get 2
        call 11
        i32.const 1
        i32.lt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i32.load offset=1423 align=1
        i32.store offset=15 align=1
        local.get 0
        i32.const 0
        i64.load offset=1416
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1408
        i64.store
        local.get 0
        i32.const 18
        call 3
        br 1 (;@1;)
      end
      local.get 2
      local.get 2
      local.get 1
      call 12
      i32.const 1120
      i32.const 32
      local.get 2
      call 7
      drop
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1024
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1028 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1036 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1044 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1052
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 2
      call 8
      drop
      block  ;; label = @2
        local.get 1
        local.get 2
        call 11
        i32.const 1
        i32.lt_s
        br_if 0 (;@2;)
        local.get 0
        i32.const 0
        i32.load offset=1391 align=1
        i32.store offset=15 align=1
        local.get 0
        i32.const 0
        i64.load offset=1384
        i64.store offset=8
        local.get 0
        i32.const 0
        i64.load offset=1376
        i64.store
        local.get 0
        i32.const 18
        call 3
        br 1 (;@1;)
      end
      local.get 2
      local.get 2
      local.get 1
      call 12
      i32.const 1120
      i32.const 32
      local.get 2
      call 7
      drop
      i32.const 0
      i32.const 1
      i32.store16 offset=1120
      i32.const 0
      i32.const 0
      i32.load offset=1056
      i32.store offset=1122 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1060 align=4
      i64.store offset=1126 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1068 align=4
      i64.store offset=1134 align=2
      i32.const 0
      i32.const 0
      i64.load offset=1076 align=4
      i64.store offset=1142 align=2
      i32.const 0
      i32.const 0
      i32.load16_u offset=1084
      i32.store16 offset=1150
      i32.const 1120
      i32.const 32
      i64.const 0
      call 5
      local.tee 2
      call 8
      drop
      local.get 2
      local.get 2
      local.get 1
      call 13
      i32.const 1120
      i32.const 32
      local.get 2
      call 7
      drop
      i32.const 1312
      i32.const 1024
      i32.const 1056
      local.get 1
      call 15
      i64.const 1
      call 14
    end
    local.get 0
    i32.const 32
    i32.add
    global.set 0)
  (table (;0;) 1 1 funcref)
  (memory (;0;) 2)
  (global (;0;) (mut i32) (i32.const 66976))
  (export "memory" (memory 0))
  (export "init" (func 16))
  (export "upgrade" (func 17))
  (export "totalSupply" (func 18))
  (export "balanceOf" (func 19))
  (export "allowance" (func 20))
  (export "transferToken" (func 21))
  (export "approve" (func 22))
  (export "transferFrom" (func 23))
  (data (;0;) (i32.const 1024) "\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00\00")
  (data (;1;) (i32.const 1280) "q4i+#\0b\9e\1f\fa9\09\89\04r!4\15\96R\b0\9c[\c4\1d\88\d6i\87y\d2(\ff\f0\99\cd\8b\deUx\14\84*1!\e8\dd\fdC:S\9b\8c\9f\14\bf1\eb\f1\08\d1.a\96\e9")
  (data (;2;) (i32.const 1344) "wrong args num\00\00negative amount\00insufficient funds\00\00\00\00\00\00\00\00\00\00\00\00\00\00allowance exceeded\00"))
