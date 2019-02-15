'use strict';

const fs = require('fs');
const path = require('path');

// Bring Fabric SDK network class
const { FileSystemWallet, Gateway } = require('fabric-network');

let walletDir = path.join(__dirname,'/_idwallet');
const wallet = new FileSystemWallet(walletDir);

const ccpPath = path.resolve(__dirname, 'connection.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const connectionProfile = JSON.parse(ccpJSON);

const channel = "mychannel";
const contractName = "mycc";
const user = "Admin";


exports.testFunction = async function(req, res, next) {
  res.send({'result': 'success'})
}

exports.changeColor = async function(req, res, next) {
  try {
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);
    const newColor = req.body.color;

    const response = await contract.submitTransaction("changeColor", newColor);

    console.log('Disconnect from Fabric gateway.');
    console.log('getMembers Complete');
    await gateway.disconnect();
    res.send({'result': 'success'});

  } catch(error) {
      console.log(`Error processing transaction. ${error}`);
      console.log(error.stack);
      res.send({'error': error.stack});
  }
}

exports.getByColor = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const color = req.body.color;

    const response = await contract.evaluateTransaction('getByColor', color);
    var cars = JSON.parse(response.toString());

    console.log('Disconnect from Fabric gateway.');
    console.log('getByColor Complete');
    await gateway.disconnect();
    res.send({'result': 'success', 'cars': cars});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}

exports.getByColorOwner = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const color = req.body.color;
    const owner = req.body.owner;

    const response = await contract.evaluateTransaction('getByColorOwner', color, owner);
    var cars = JSON.parse(response.toString());

    console.log('Disconnect from Fabric gateway.');
    console.log('getByColorOwner Complete');
    await gateway.disconnect();
    res.send({'result': 'success', 'orders': cars});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}

exports.getCarHistory = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const chassis = req.body.chassis;

    const response = await contract.evaluateTransaction('getCarHistory', chassis);
    var cars = JSON.parse(response.toString());

    console.log('Disconnect from Fabric gateway.');
    console.log('getCarHistory Complete');
    await gateway.disconnect();
    res.send({'result': 'success', 'orders': cars});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }

}

exports.changeOwner = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const chassis = req.body.chassis;
    const buyer = req.body.buyer;
    const payTheDmg = req.body.payTheDmg;

    const response = await contract.submitTransaction('getCarHistory', chassis, buyer, payTheDmg);

    console.log('Disconnect from Fabric gateway.');
    console.log('changeOwner Complete');
    await gateway.disconnect();
    res.send({'result': 'success'});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}

exports.noteDamage = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const chassis = req.body.chassis;
    const description = req.body.description;
    const repairPrice = req.body.repairPrice;

    const response = await contract.submitTransaction('noteDamage', chassis, description, repairPrice);

    console.log('Disconnect from Fabric gateway.');
    console.log('noteDamage Complete');
    await gateway.disconnect();
    res.send({'result': 'success'});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}

exports.repairDamage = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const chassis = req.body.chassis;

    const response = await contract.submitTransaction('repairDamage', chassis);
    console.log('Disconnect from Fabric gateway.');
    console.log('repairDamage Complete');
    await gateway.disconnect();
    res.send({'result': 'success'});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}

exports.initLedger = async function(req, res, next) {
  try {
    
    const gateway = new Gateway();
    await gateway.connect(connectionProfile, { wallet, identity: user, discovery: { enabled: false } });
    const network = await gateway.getNetwork(channel);

    const contract = await network.getContract(contractName);

    const response = await contract.submitTransaction('initLedger', '');
    console.log('Disconnect from Fabric gateway.');
    console.log('initLedger Complete');
    await gateway.disconnect();
    res.send({'result': 'success'});
  } catch (error) {
    console.log(`Error processing transaction. ${error}`);
    console.log(error.stack);
    res.send({'error': error.stack});
  }
}