#set -x
#OIFS=$IFS
#IFS=''
export PREFIX=/tank/media/games/PSX
#ls -C "${PREFIX}" | tail -n 2 > files.txt
ls -C "${PREFIX}" > files.txt
while read -r filename ; do    
  echo "FILENAME: $filename"
  SOURCE="${PREFIX}/${filename}"
  ROOTFN=`echo ${filename} | cut -f1 -d.`
  TEMP="${ROOTFN}.$$"
  DEST="${ROOTFN}"
  mkdir -p "${DEST}" "${TEMP}"
  unzip "${SOURCE}" -d "${TEMP}"
  set -x
  echo "`ls "$TEMP/" | grep cue`" > cuefilename
  cuefn=`cat cuefilename`
  chdman createcd -i "${TEMP}/${cuefn}" -o "${DEST}/${ROOTFN}.chd"
  rm -fr "$TEMP"
  set +x
done < files.txt
exit 0;

