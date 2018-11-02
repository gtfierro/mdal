from setuptools import setup, find_packages

print(find_packages())
setup(name='dataclient',
      version='0.1.10',
      description='MDAL client v2',
      packages=find_packages(),
      data_files=[('dataclient', ['dataclient/data.capnp'])],
      include_package_data=True,
      #install_requires=[
      #  'delorean==0.6.0',
      #  'msgpack-python==0.4.2',
      #  'requests>=2.12.2',
      #  'python-dateutil>=2.4.2',
      #  'pandas>=0.20.1',
      #  'pycapnp>=0.6.3',
      #],
      zip_safe=False)

