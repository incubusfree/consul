// rootURL in production equals `{{.ContentPath}}` and therefore is replaced
// with the value of -ui-content-path. During development rootURL uses the
// value as set in environment.js
module.exports = ({ appName, environment, rootURL, config }) => `
  <noscript>
      <div style="margin: 0 auto;">
          <h2>JavaScript Required</h2>
          <p>Please enable JavaScript in your web browser to use Consul UI.</p>
      </div>
  </noscript>

  <div class="brand-loader">
    <svg width="${
      config.CONSUL_BINARY_TYPE !== 'oss' && config.CONSUL_BINARY_TYPE !== '' ? `394` : `198`
    }" height="53" xmlns="http://www.w3.org/2000/svg" fill="#919FA8">
      <path d="M32.7240001,0.866235051 C28.6239001,-0.218137949 24.3210001,-0.285465949 20.1890001,0.670096051 C16.0569001,1.62566005 12.2205001,3.57523005 9.01276015,6.34960005 C5.80499015,9.12397005 3.32280015,12.6393001 1.78161015,16.5905001 C0.240433148,20.5416001 -0.313157852,24.8092001 0.168892148,29.0228001 C0.650943148,33.2364001 2.15407015,37.2687001 4.54780015,40.7697001 C6.94153015,44.2707001 10.1535001,47.1346001 13.9050001,49.1128001 C17.6565001,51.0910001 21.8341001,52.1238001 26.0752001,52.1214409 C32.6125001,52.1281001 38.9121001,49.6698001 43.7170001,45.2370001 L37.5547001,38.7957001 C35.0952001,41.0133001 32.0454001,42.4701001 28.7748001,42.9898001 C25.5042001,43.5096001 22.1530001,43.0698001 19.1273001,41.7239001 C16.1015001,40.3779001 13.5308001,38.1835001 11.7267001,35.4064001 C9.92260015,32.6294001 8.96239015,29.3888001 8.96239015,26.0771001 C8.96239015,22.7655001 9.92260015,19.5249001 11.7267001,16.7478001 C13.5308001,13.9707001 16.1015001,11.7763001 19.1273001,10.4304001 C22.1530001,9.08444005 25.5042001,8.64470005 28.7748001,9.16441005 C32.0454001,9.68412005 35.0952001,11.1410001 37.5547001,13.3586001 L43.7170001,6.89263005 C40.5976001,4.01926005 36.8241001,1.95061005 32.7240001,0.866235051 Z M46.6320001,34.8572001 C46.2182001,34.9395001 45.8380001,35.1427001 45.5397001,35.4410001 C45.2413001,35.7394001 45.0381001,36.1195001 44.9558001,36.5334001 C44.8735001,36.9472001 44.9157001,37.3762001 45.0772001,37.7660001 C45.2387001,38.1559001 45.5121001,38.4891001 45.8630001,38.7235001 C46.2138001,38.9579001 46.6263001,39.0830001 47.0482001,39.0830001 C47.6141001,39.0830001 48.1567001,38.8583001 48.5568001,38.4582001 C48.9569001,38.0581001 49.1817001,37.5154001 49.1817001,36.9496001 C49.1817001,36.5276001 49.0565001,36.1152001 48.8221001,35.7643001 C48.5877001,35.4135001 48.2545001,35.1400001 47.8647001,34.9786001 C47.4748001,34.8171001 47.0459001,34.7748001 46.6320001,34.8572001 Z M49.0856001,27.5622001 C48.6718001,27.6446001 48.2916001,27.8477001 47.9933001,28.1461001 C47.6949001,28.4445001 47.4917001,28.8246001 47.4094001,29.2385001 C47.3271001,29.6523001 47.3693001,30.0813001 47.5308001,30.4711001 C47.6923001,30.8609001 47.9657001,31.1941001 48.3166001,31.4286001 C48.6674001,31.6630001 49.0799001,31.7881001 49.5018001,31.7881001 C50.0670001,31.7859001 50.6084001,31.5605001 51.0080001,31.1609001 C51.4076001,30.7612001 51.6331001,30.2198001 51.6353001,29.6547001 C51.6353001,29.2327001 51.5102001,28.8202001 51.2757001,28.4694001 C51.0413001,28.1186001 50.7081001,27.8451001 50.3183001,27.6836001 C49.9284001,27.5222001 49.4995001,27.4799001 49.0856001,27.5622001 Z M28.0728001,20.8457001 C27.0412001,20.4185001 25.9061001,20.3067001 24.8110001,20.5245001 C23.7159001,20.7423001 22.7100001,21.2800001 21.9205001,22.0695001 C21.1309001,22.8590001 20.5933001,23.8650001 20.3754001,24.9600001 C20.1576001,26.0551001 20.2694001,27.1902001 20.6967001,28.2218001 C21.1240001,29.2534001 21.8476001,30.1351001 22.7760001,30.7554001 C23.7043001,31.3757001 24.7958001,31.7068001 25.9124001,31.7068001 C27.4096001,31.7068001 28.8455001,31.1120001 29.9043001,30.0533001 C30.9630001,28.9946001 31.5578001,27.5587001 31.5578001,26.0614001 C31.5578001,24.9449001 31.2267001,23.8534001 30.6063001,22.9250001 C29.9860001,21.9966001 29.1043001,21.2730001 28.0728001,20.8457001 Z M43.9670001,27.4378001 C43.5772001,27.2763001 43.1482001,27.2341001 42.7344001,27.3164001 C42.3205001,27.3987001 41.9404001,27.6019001 41.6420001,27.9003001 C41.3437001,28.1986001 41.1405001,28.5788001 41.0581001,28.9926001 C40.9758001,29.4065001 41.0181001,29.8354001 41.1796001,30.2253001 C41.3410001,30.6151001 41.6145001,30.9483001 41.9653001,31.1827001 C42.3162001,31.4171001 42.7286001,31.5423001 43.1506001,31.5423001 C43.7164001,31.5423001 44.2591001,31.3175001 44.6592001,30.9174001 C45.0592001,30.5173001 45.2840001,29.9747001 45.2840001,29.4088001 C45.2840001,28.9869001 45.1589001,28.5744001 44.9245001,28.2236001 C44.6901001,27.8727001 44.3568001,27.5993001 43.9670001,27.4378001 Z M43.9670001,20.7503001 C43.5772001,20.5888001 43.1482001,20.5466001 42.7344001,20.6289001 C42.3205001,20.7112001 41.9404001,20.9144001 41.6420001,21.2128001 C41.3437001,21.5111001 41.1405001,21.8913001 41.0581001,22.3051001 C40.9758001,22.7190001 41.0181001,23.1479001 41.1796001,23.5378001 C41.3410001,23.9276001 41.6145001,24.2608001 41.9653001,24.4952001 C42.3162001,24.7296001 42.7286001,24.8548001 43.1506001,24.8548001 C43.7164001,24.8548001 44.2591001,24.6300001 44.6592001,24.2299001 C45.0592001,23.8298001 45.2840001,23.2871001 45.2840001,22.7213001 C45.2840001,22.2994001 45.1589001,21.8869001 44.9245001,21.5360001 C44.6901001,21.1852001 44.3568001,20.9118001 43.9670001,20.7503001 Z M49.0856001,20.3825001 C48.6718001,20.4649001 48.2916001,20.6681001 47.9933001,20.9664001 C47.6949001,21.2648001 47.4917001,21.6449001 47.4094001,22.0588001 C47.3271001,22.4726001 47.3693001,22.9016001 47.5308001,23.2914001 C47.6923001,23.6813001 47.9657001,24.0144001 48.3166001,24.2489001 C48.6674001,24.4833001 49.0799001,24.6084001 49.5018001,24.6084001 C50.0670001,24.6063001 50.6084001,24.3808001 51.0080001,23.9812001 C51.4076001,23.5815001 51.6331001,23.0401001 51.6353001,22.4750001 C51.6353001,22.0530001 51.5102001,21.6406001 51.2757001,21.2897001 C51.0413001,20.9389001 50.7081001,20.6654001 50.3183001,20.5040001 C49.9284001,20.3425001 49.4995001,20.3002001 49.0856001,20.3825001 Z M46.7554001,13.2026001 C46.3416001,13.2849001 45.9614001,13.4881001 45.6630001,13.7865001 C45.3647001,14.0849001 45.1615001,14.4650001 45.0792001,14.8788001 C44.9969001,15.2927001 45.0391001,15.7217001 45.2006001,16.1115001 C45.3621001,16.5013001 45.6355001,16.8345001 45.9863001,17.0689001 C46.3372001,17.3034001 46.7497001,17.4285001 47.1716001,17.4285001 C47.7374001,17.4285001 48.2801001,17.2037001 48.6802001,16.8036001 C49.0803001,16.4035001 49.3050001,15.8609001 49.3050001,15.2951001 C49.3050001,14.8731001 49.1799001,14.4606001 48.9455001,14.1098001 C48.7111001,13.7589001 48.3779001,13.4855001 47.9880001,13.3240001 C47.5982001,13.1625001 47.1692001,13.1203001 46.7554001,13.2026001 Z" fill-rule="nonzero"></path>
      <path d="M83.5385001,9.02612084 C75.3002001,9.02612084 71.7718001,12.5545001 71.7718001,18.6102001 L71.7718001,33.5278001 L71.7744126,33.809806 C71.8842215,39.6928981 75.4612111,43.1118103 83.5385001,43.1118103 C86.5802001,43.1131001 89.6109001,42.7466001 92.5646001,42.0205001 L91.8671001,36.6049001 L90.9760579,36.7631811 C88.5964705,37.1629803 86.1899224,37.3844223 83.7765001,37.4254001 C79.4194001,37.4254001 78.0326001,35.9320001 78.0326001,32.4118001 L78.0326001,19.7261001 L78.0346281,19.4988781 C78.0956946,16.133828 79.5462067,14.7125001 83.7765001,14.7125001 C86.4916001,14.7587001 89.1980001,15.0332001 91.8671001,15.5331001 L92.5646001,10.1175001 L91.8246092,9.94345672 C89.1057071,9.33281156 86.3267251,9.02500229 83.5385001,9.02612084 Z M172.149,18.4131001 L166.094,18.4131001 L166.09588,36.2248122 C166.154955,40.3975255 167.61375,43.1117001 171.55,43.1117001 C174.919,42.9517001 178.218,42.0880001 181.233,40.5762001 L181.832,42.6112001 L186.443,42.6112001 L186.443,18.4131001 L180.388,18.4131001 L180.388,35.1934001 C178.188,36.3339001 175.481,37.2283001 174.086,37.2283001 C172.691,37.2283001 172.149,36.5801001 172.149,35.2918001 L172.149,18.4131001 Z M105.939,17.9127001 C98.2719471,17.9127001 95.7845671,21.8519543 95.4516942,26.3358062 L95.4257941,26.7784774 C95.4225999,26.8525088 95.4199581,26.9266566 95.4178553,27.0009059 L95.4116001,27.4475001 L95.4116001,33.5853001 L95.4178331,34.0318054 C95.5519456,38.7818866 97.886685,43.0872001 105.931,43.0872001 C113.716697,43.0872001 116.15821,39.0467642 116.432186,34.4757046 L116.45204,34.0318054 C116.456473,33.8833653 116.458758,33.734491 116.459,33.5853001 L116.459,27.4475001 L116.457455,27.2221358 C116.453317,26.9220505 116.440796,26.6236441 116.419035,26.3278463 L116.379357,25.8862225 C115.91894,21.5651129 113.355121,17.9127001 105.939,17.9127001 Z M154.345,17.8876515 C147.453,17.8876515 145.319,20.0214001 145.319,24.8873001 L145.319694,25.1343997 L145.325703,25.6107983 L145.338905,26.064173 C145.341773,26.1378641 145.344992,26.2106314 145.348588,26.2824927 L145.374889,26.7029295 C145.380095,26.7712375 145.385729,26.838675 145.391816,26.9052596 L145.433992,27.2946761 C145.714183,29.5082333 146.613236,30.7206123 149.232713,31.693068 L149.698825,31.8575665 C150.021076,31.9658547 150.36662,32.0715774 150.737101,32.1758709 L151.311731,32.3313812 C151.509646,32.3829554 151.714,32.4343143 151.925,32.4856001 L152.205551,32.5543061 L152.728976,32.6899356 L153.204098,32.8237311 L153.633238,32.9563441 C155.53221,33.5734587 156.004908,34.1732248 156.112605,35.0535762 L156.130482,35.2466262 L156.139507,35.448917 L156.142,35.6611001 L156.137247,35.9859786 L156.121298,36.2838969 C156.024263,37.5177444 155.540462,38.0172149 153.741624,38.1073495 L153.302742,38.1210314 L153.065,38.1227001 C150.631,38.0987001 148.21,37.7482001 145.869,37.0807001 L145.049,41.6922001 L145.672496,41.887484 C148.174444,42.639635 150.769923,43.0436231 153.385,43.0871001 C159.627887,43.0871001 161.583469,40.9824692 162.030289,37.4548504 L162.074576,37.049455 C162.087289,36.9123213 162.098004,36.7731979 162.106868,36.6321214 L162.128062,36.2030694 L162.139051,35.7625187 L162.141,35.5380001 C162.141,35.4566181 162.140828,35.3763299 162.14046,35.2971136 L162.131203,34.6125174 L162.117224,34.1865271 L162.095649,33.7836378 L162.065324,33.4027996 L162.025093,33.0429627 L161.973799,32.7030773 C161.659145,30.8866498 160.790109,29.9278873 158.501441,29.0408119 L158.069484,28.8801405 L157.605084,28.7199991 C157.524916,28.6932947 157.443348,28.6665687 157.360357,28.6397991 L156.845127,28.4784845 L156.294565,28.3150754 L155.707516,28.148522 L155.082823,27.9777746 L154.035614,27.7021396 L153.423677,27.5325226 L153.071612,27.4262327 C153.016479,27.4088193 152.963082,27.3915263 152.911366,27.3743086 L152.620815,27.2715428 C151.671458,26.912485 151.415595,26.5466416 151.348761,25.7543883 L151.334373,25.5160648 L151.327658,25.2523603 L151.327351,24.8244501 C151.355827,23.4390475 151.851313,22.8769001 154.403,22.8769001 C156.636,22.9360001 158.861,23.1692001 161.057,23.5744001 L161.591,18.7085001 L160.876597,18.5511522 C158.72872,18.1040608 156.5401,17.8816774 154.345,17.8876515 Z M197.71,7.71350005 L191.654,8.53405005 L191.654,42.6116001 L197.71,42.6116001 L197.71,7.71350005 Z M135.455,17.9211001 C132.086,18.0823001 128.788,18.9459001 125.772,20.4566001 L125.189,18.4135001 L120.57,18.4135001 L120.57,42.6115001 L126.625,42.6115001 L126.625,25.8066001 C128.833,24.6661001 131.549,23.7717001 132.936,23.7717001 C134.322,23.7717001 134.872,24.4199001 134.872,25.7082001 L134.872,42.6115001 L140.919,42.6115001 L140.919,25.0681001 C140.919,20.7520001 139.475,17.9211001 135.455,17.9211001 Z M105.931,23.0740001 C109.156,23.0740001 110.395,24.5592001 110.395,27.2506001 L110.395,33.7494001 L110.392134,33.9740961 C110.325067,36.5604698 109.074195,37.9178001 105.931,37.9178001 C102.698,37.9178001 101.459,36.4818001 101.459,33.7494001 L101.459,27.2506001 L101.461884,27.0258853 C101.529372,24.4390811 102.787806,23.0740001 105.931,23.0740001 Z" fill-rule="nonzero"></path>
      ${
        config.CONSUL_BINARY_TYPE !== 'oss' && config.CONSUL_BINARY_TYPE !== ''
          ? `
        <path d="M322.099,18.0445001 C319.225,18.0223001 316.427,18.9609001 314.148,20.7112001 L314.016,20.8179001 L313.68,18.5368001 L310.332,18.5368001 L310.332,53.0000001 L314.312,52.4338001 L314.312,42.3164001 L314.435,42.3164001 C316.705,42.7693001 319.012,43.0165001 321.327,43.0549001 C326.554,43.0549001 329.098,40.5029001 329.098,35.2432001 L329.098,25.3802001 C329.073,20.4569001 326.809,18.0445001 322.099,18.0445001 Z M264.971,11.9722001 L260.991,12.5466001 L260.991,18.5284001 L256.708,18.5284001 L256.708,21.8106001 L260.991,21.8106001 L260.991,37.6883001 L260.99344,37.9365729 C261.066744,41.6122056 262.7975,43.1124033 266.915,43.1124033 C268.591,43.1170001 270.255,42.8396001 271.839,42.2915001 L271.363,39.1817001 L270.896229,39.3066643 C269.803094,39.5806719 268.682875,39.7315001 267.555,39.7560001 C265.526625,39.7560001 265.081547,38.9674128 264.991981,37.7056542 L264.97743,37.4176027 L264.97159,37.1147428 L264.971,21.8188001 L271.494,21.8188001 L271.83,18.5366001 L264.971,18.5366001 L264.971,11.9722001 Z M283.556,18.0770001 C277.312,18.0770001 274.144,21.0884001 274.144,27.0374001 L274.144,34.3075001 C274.144,40.3140001 277.164,43.1124894 283.655,43.1124894 C286.526,43.1192001 289.38,42.6620001 292.106,41.7581001 L291.589,38.6154001 C289.116,39.3030001 286.566,39.6779001 283.999,39.7314001 C279.785843,39.7314001 278.500803,38.4772648 278.201322,35.860808 L278.165734,35.4868687 L278.141767,35.0951811 C278.138675,35.0284172 278.136019,34.9609111 278.133774,34.8926614 L278.125037,34.474229 L278.124,32.0756001 L292.582,32.0756001 L292.582,27.1031001 C292.582,21.0064001 289.636,18.0770001 283.556,18.0770001 Z M384.631,18.0768001 C378.412,18.0440001 375.22,21.0554001 375.22,27.0208001 L375.22,34.2909001 C375.22,40.2973001 378.239,43.0955988 384.73,43.0955988 C387.599,43.1033001 390.45,42.6460001 393.173,41.7415001 L392.665,38.5988001 C390.188,39.2815001 387.635,39.6509001 385.066,39.6983001 C380.852843,39.6983001 379.567803,38.4442359 379.268322,35.8278014 L379.232734,35.4538649 L379.208767,35.0621794 C379.205675,34.9954158 379.203019,34.9279099 379.200774,34.8596604 L379.192037,34.4412289 L379.191,32.0754001 L393.657,32.0754001 L393.657,27.1029001 C393.657,21.0062001 390.712,18.0768001 384.631,18.0768001 Z M364.634,18.0441001 C363.881125,18.0441001 363.18736,18.0712813 362.54969,18.1279834 L362.016783,18.1838695 C357.948857,18.6791301 356.371,20.5353768 356.371,24.4608001 L356.371522,24.7155013 L356.376145,25.2052033 L356.386527,25.669464 L356.403852,26.1092746 C356.407384,26.1805939 356.411254,26.2509357 356.415488,26.3203208 L356.445451,26.7253144 L356.485319,27.1083357 C356.756619,29.3425283 357.626845,30.4437319 360.247859,31.3753061 L360.701103,31.529163 C360.779411,31.5545991 360.85912,31.5799457 360.940253,31.6052232 L361.444353,31.7562266 L361.983836,31.9065664 L362.55989,32.0572338 L363.430663,32.2724269 L364.440153,32.5299129 L364.884369,32.6506971 L365.29049,32.7679922 L365.660213,32.8831607 L365.99523,32.9975651 C367.26815,33.4554713 367.748817,33.9277406 367.925217,34.806783 L367.963261,35.0352452 C367.974017,35.1143754 367.982943,35.1965576 367.990321,35.2820187 L368.008092,35.5484662 L368.018269,35.8359502 L368.023,36.3096001 C368.023,36.3683432 368.022674,36.4261667 368.021989,36.4830819 L368.013333,36.8137655 C368.008847,36.9204214 368.002676,37.0235359 367.994568,37.1232009 L367.964177,37.4119383 C367.774513,38.8512264 367.058626,39.4837671 364.875404,39.6510671 L364.43427,39.67773 L363.954974,39.6933243 C363.78868,39.6967387 363.615773,39.6984001 363.436,39.6984001 C361.126,39.6638001 358.83,39.3385001 356.601,38.7302001 L356.051,41.7908001 L356.619468,41.9710684 C358.900888,42.6645722 361.270923,43.0269154 363.658,43.0463001 C369.59355,43.0463001 371.402903,41.3625861 371.812159,38.0405419 L371.854011,37.6421573 C371.859965,37.574501 371.865421,37.5062155 371.870401,37.4373012 L371.894725,37.0162715 L371.908596,36.5801656 C371.911587,36.4322862 371.913,36.2818967 371.913,36.1290001 L371.914417,35.5317322 C371.901583,33.4289389 371.677,32.2649251 370.797,31.3698001 C370.053077,30.6022731 368.787947,30.0494771 366.870096,29.4840145 L366.242608,29.3047611 C366.13436,29.2747269 366.024265,29.2445914 365.912304,29.2143213 L365.218,29.0308209 L364.216102,28.7784328 L363.495981,28.593015 L363.068145,28.4733265 L362.67987,28.3551624 C361.018765,27.8247783 360.501056,27.2986662 360.340522,26.2094051 L360.310407,25.9578465 C360.306262,25.9142982 360.302526,25.8699197 360.29916,25.8246823 L360.283089,25.5427193 L360.273984,25.2387571 L360.269927,24.911412 L360.270221,24.3885398 L360.280627,24.0635689 C360.366727,22.3885604 360.966747,21.6370879 363.248047,21.4645754 L363.695778,21.4389299 L364.184625,21.426349 L364.445,21.4248001 C366.684,21.4608001 368.916,21.6859001 371.117,22.0976001 L371.396,18.8646001 L370.730951,18.7059457 C368.73071,18.2553391 366.686,18.0331201 364.634,18.0441001 Z M351.301,18.5363001 L347.321,18.5363001 L347.321,42.6112001 L351.301,42.6112001 L351.301,18.5363001 Z M307.335,18.0850001 L306.70097,18.3638937 C304.598769,19.3169298 302.610091,20.5031364 300.771,21.9005001 L300.623,22.0236001 L300.369,18.5363001 L296.931,18.5363001 L296.931,42.6112001 L300.91,42.6112001 L300.91,25.9048001 L301.641825,25.3925123 C303.604371,24.0427531 305.654445,22.8240667 307.778,21.7446001 L307.335,18.0850001 Z M344.318,18.0850001 L343.683947,18.3638937 C341.581595,19.3169298 339.592091,20.5031364 337.753,21.9005001 L337.606,22.0236001 L337.351,18.5363001 L333.946,18.5363001 L333.946,42.6112001 L337.926,42.6112001 L337.926,25.9048001 L337.967,25.9048001 L338.701162,25.3884311 C340.669963,24.0279284 342.726556,22.7996223 344.859,21.7118001 L344.318,18.0850001 Z M230.384,9.62500005 L211.109,9.62500005 L211.109,42.6112001 L230.466,42.6112001 L230.466,38.9597001 L215.146,38.9597001 L215.146,27.4720001 L229.293,27.4720001 L229.293,23.8698001 L215.146,23.8698001 L215.146,13.2600001 L230.384,13.2600001 L230.384,9.62500005 Z M248.763,18.0441001 C245.899,18.0441001 241.706,19.3323001 239.047,20.6124001 L238.924,20.6698001 L238.522,18.5282001 L235.322,18.5282001 L235.322,42.5704001 L239.302,42.5704001 L239.302,24.2885001 L239.359,24.2885001 C241.919,22.9674001 245.661,21.8268001 247.524,21.8268001 C249.165,21.8268001 249.985,22.5735001 249.985,24.1736001 L249.985,42.5868001 L253.965,42.5868001 L253.965,24.1161001 C253.932,20.0380001 252.25,18.0523001 248.763,18.0441001 Z M321.229,21.5564001 C323.526,21.5564001 325.061,22.2046001 325.061,25.3966001 L325.094,35.2760001 C325.094,38.3121001 323.887,39.6085001 321.057,39.6085001 C318.81,39.5533001 316.572,39.3035001 314.369,38.8618001 L314.287,38.8618001 L314.287,24.4694001 C316.198,22.7311001 318.649,21.7027001 321.229,21.5564001 Z M283.581,21.3264001 C287.372,21.3264001 288.758,22.8855001 288.758,26.7010001 L288.758,28.7934001 L278.149,28.7934001 L278.149,26.7010001 C278.149,22.9839001 279.79,21.3264001 283.581,21.3264001 Z M384.648,21.3262001 C388.431,21.3262001 389.834,22.8852001 389.834,26.7008001 L389.834,28.7932001 L379.224,28.7932001 L379.224,26.7008001 C379.224,22.9837001 380.865,21.3262001 384.648,21.3262001 Z M351.301,8.63220005 L347.321,8.63220005 L347.321,14.4499001 L351.301,14.4499001 L351.301,8.63220005 Z" fill-rule="nonzero"></path>
      `
          : ``
      }
    </svg>
  </div>
  <script type="application/json" data-consul-ui-config>
${environment === 'production' ? `{{jsonEncode .}}` : JSON.stringify(config.operatorConfig)}
  </script>
  <script type="application/json" data-consul-ui-fs>
  {
    "text-encoding/encoding-indexes.js": "${rootURL}assets/encoding-indexes.js",
    "text-encoding/encoding.js": "${rootURL}assets/encoding-indexes.js",
    "css.escape/css.escape.js": "${rootURL}assets/css.escape.js",
    "codemirror/mode/javascript/javascript.js": "${rootURL}assets/codemirror/mode/javascript/javascript.js",
    "codemirror/mode/ruby/ruby.js": "${rootURL}assets/codemirror/mode/ruby/ruby.js",
    "codemirror/mode/yaml/yaml.js": "${rootURL}assets/codemirror/mode/yaml/yaml.js"
  }
  </script>
  <script src="${rootURL}assets/init.js"></script>
  <script src="${rootURL}assets/vendor.js"></script>
  ${environment === 'test' ? `<script src="${rootURL}assets/test-support.js"></script>` : ``}
  <script src="${rootURL}assets/metrics-providers/consul.js"></script>
  <script src="${rootURL}assets/metrics-providers/prometheus.js"></script>
  ${
    environment === 'production'
      ? `{{ range .ExtraScripts }} <script src="{{.}}"></script> {{ end }}`
      : ``
  }
  <script src="${rootURL}assets/${appName}.js"></script>
  ${environment === 'test' ? `<script src="${rootURL}assets/tests.js"></script>` : ``}
`;
