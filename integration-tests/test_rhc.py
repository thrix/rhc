def test_version(rhc):
    proc = rhc.run("--version")
    assert b"rhc version" in proc.stdout


def test_rhc_register(external_candlepin, test_config, rhc):
    assert not rhc.is_registered
    rhc.connect(
        username=test_config.get("candlepin", "username"),
        password=test_config.get("candlepin", "password"),
    )
    assert rhc.is_registered
    rhc.disconnect()
    assert not rhc.is_registered
