# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.define "db1" do |db1|
    db1.vm.box = "centos/7"
    db1.vm.network "private_network", ip: "192.168.33.51"
    db1.vm.provider "virtualbox" do |vb|
      disk_file = "db1sdb.vdi"
      unless File.exists?(disk_file)
        vb.customize [
          'createmedium', 'disk',
          '--filename', disk_file,
          '--format', 'VDI',
          '--size', 20 * 1024]
      end
      vb.customize [
        'storageattach', :id,
        '--storagectl', 'IDE',
        '--port', 1,
        '--device', 0,
        '--type', 'hdd',
        '--medium', disk_file]
      vb.memory = "2048"
      vb.cpus = 2
    end
  end

  config.vm.define "db2" do |db2|
    db2.vm.box = "centos/7"
    db2.vm.network "private_network", ip: "192.168.33.52"
    db2.vm.provider "virtualbox" do |vb|
      disk_file = "db2sdb.vdi"
      unless File.exists?(disk_file)
        vb.customize [
          'createmedium', 'disk',
          '--filename', disk_file,
          '--format', 'VDI',
          '--size', 20 * 1024]
      end
      vb.customize [
        'storageattach', :id,
        '--storagectl', 'IDE',
        '--port', 1,
        '--device', 0,
        '--type', 'hdd',
        '--medium', disk_file]
      vb.memory = "2048"
      vb.cpus = 2
    end
  end

end
